package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/mdlayher/unifi"
)

type wlitem struct {
	Alias      string         `json:"alias"`
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
	template *template.Template
}

type view struct {
	Self wlitem
	List wlitems
}

func NewWhitelist(c *unifi.Client, jsonfile string) (wl *whitelist) {
	return &whitelist{
		List:     make(wlitems),
		stations: make(stations),
		client:   c,
		jsonfile: jsonfile,
		template: template.Must(template.New("index").Parse(string(MustAsset("index.tmpl")))),
	}
}

func newWlitem(alias string, s *unifi.Station) (wi *wlitem) {
	return (&wlitem{
		Alias:  alias,
		Online: true,
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

func (wl *whitelist) view(r *http.Request) view {
	self := wl.self(r)
	alias := self.Hostname
	if wi, ok := wl.List[self.MAC.String()]; ok {
		alias = wi.Alias
	}
	return view{*newWlitem(alias, self), wl.List}
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

func (wl *whitelist) OnChange() {
	fmt.Println("change detected")
}

func (wl *whitelist) Update() (err error) {
	var ss []*unifi.Station

	ss, err = wl.client.Stations(config.site)
	if err == nil {
		changed := false
		wl.stations = make(stations)
		for _, s := range ss {
			wl.stations[s.MAC.String()] = s
		}
		for mac := range wl.List {
			s, ok := wl.stations[mac]
			changed = changed || wl.List[mac].Online != ok
			if ok {
				wl.List[mac].syncFrom(s)
			}
			wl.List[mac].Online = ok
		}
		if changed {
			wl.OnChange()
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
	wl.template.Execute(w, wl.view(r))
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
