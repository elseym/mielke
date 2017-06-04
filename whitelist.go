package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mdlayher/unifi"
)

type wlitem struct {
	Alias      string         `json:"alias"`
	AvatarURL  string         `json:"avatarURL"`
	Hostname   string         `json:"hostname"`
	Online     bool           `json:"online"`
	Associated time.Time      `json:"associated"`
	LastSeen   time.Time      `json:"lastSeen"`
	AP         string         `json:"ap"`
	Info       *unifi.Station `json:"-"`
}

type wlitems map[string]*wlitem
type stations map[string]*unifi.Station

type whitelist struct {
	List     wlitems
	stations stations
	client   *unifi.Client
	jsonfile string
}

type view struct {
	Self wlitem    `json:"self"`
	List []*wlitem `json:"list"`
}

func (v view) Len() int           { return len(v.List) }
func (v view) Swap(i, j int)      { v.List[i], v.List[j] = v.List[j], v.List[i] }
func (v view) Less(i, j int) bool { return v.List[i].Hostname < v.List[j].Hostname }

func NewWhitelist(c *unifi.Client, jsonfile string) (wl *whitelist) {
	return &whitelist{
		List:     make(wlitems),
		stations: make(stations),
		client:   c,
		jsonfile: jsonfile,
	}
}

func gravatarURL(email string) (url string) {
	h := md5.New()
	io.WriteString(h, email)
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x.jpg", h.Sum(nil))
}

func newWlitem(alias string, s *unifi.Station) (wi *wlitem) {
	return (&wlitem{
		Alias:     alias,
		AvatarURL: gravatarURL(alias),
		Online:    true,
	}).syncFrom(s)
}

func (wi *wlitem) syncFrom(s *unifi.Station) *wlitem {
	wi.AP = s.APMAC.String()
	wi.Hostname = s.Hostname
	wi.Associated = s.AssociationTime
	wi.LastSeen = s.LastSeen
	wi.Info = s
	return wi
}

func (wl *whitelist) view(r *http.Request) (v view) {
	self := wl.self(r)
	alias := self.Hostname
	if wi, ok := wl.List[self.MAC.String()]; ok {
		alias = wi.Alias
	}

	list := make([]*wlitem, 0, len(wl.List))
	for _, i := range wl.List {
		list = append(list, i)
	}

	v = view{*newWlitem(alias, self), list}
	sort.Sort(v)

	return
}

func (wl *whitelist) Load() (w *whitelist, err error) {
	var data []byte
	var list wlitems

	if data, err = ioutil.ReadFile(wl.jsonfile); err != nil {
		if _, er2 := os.Stat(wl.jsonfile); os.IsNotExist(er2) {
			err = wl.Save()
		}
	} else {
		if err = json.Unmarshal(data, &list); err == nil {
			wl.List = list
		}
	}

	return wl, err
}

func (wl *whitelist) Save() (err error) {
	var data []byte

	if data, err = json.Marshal(wl.List); err == nil {
		err = ioutil.WriteFile(wl.jsonfile, data, 0644)
	}

	if err != nil {
		fmt.Println("Error saving whitelist: " + err.Error())
	}

	return
}

func (wl *whitelist) OnChange(changelist wlitems) {
	fmt.Println("change detected", changelist)
}

func (wl *whitelist) Update() (err error) {
	var ss []*unifi.Station
	var changelist = make(wlitems)

	ss, err = wl.client.Stations(config.site)
	if err == nil {
		wl.stations = make(stations)
		for _, s := range ss {
			wl.stations[s.MAC.String()] = s
		}
		for mac := range wl.List {
			s, ok := wl.stations[mac]
			if ok {
				wl.List[mac].syncFrom(s)
			}
			if wl.List[mac].Online != ok {
				changelist[mac] = wl.List[mac]
			}
			wl.List[mac].Online = ok
		}
		if len(changelist) > 0 {
			wl.OnChange(changelist)
		}
	}

	return
}

func (wl *whitelist) UpdateLoop(freq time.Duration) {
	for {
		if err := wl.Update(); err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(freq)
	}
}

func (wl *whitelist) self(r *http.Request) (station *unifi.Station) {
	ip := net.ParseIP(r.Header.Get("X-FORWARDED-FOR"))
	if ip == nil {
		ip = net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
	}

	for _, s := range wl.stations {
		if s.IP.String() == ip.String() {
			station = s
			break
		}
	}

	return
}

type ssm map[string]string

func (wl *whitelist) Router() http.Handler {
	r := router{make([]route, 0)}
	r.add("/", "GET", ssm{"Accept": "application/json"}, wl.JSONListHandler)
	r.add("/", "GET", ssm{"Accept": "text/html"}, wl.HTMLListHandler)
	r.add("/", "PUT", ssm{"Content-Type": "application/json"}, wl.AddHandler)
	r.add("/", "DELETE", ssm{}, wl.RemoveHandler)

	return r
}

func (wl *whitelist) JSONListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if data, err := json.Marshal(wl.view(r)); err == nil {
		w.Write(data)
	} else {
		w.WriteHeader(500)
	}
}

func (wl *whitelist) HTMLListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write(MustAsset("mielke.html"))
}

func (wl *whitelist) AddHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Alias string `json:"alias"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}

	self := wl.self(r)
	if data.Alias == "" {
		data.Alias = self.Hostname
	}

	wl.List[self.MAC.String()] = newWlitem(data.Alias, self)
	wl.Save()
}

func (wl *whitelist) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	delete(wl.List, wl.self(r).MAC.String())
	wl.Save()
}
