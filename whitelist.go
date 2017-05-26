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

	"github.com/mdlayher/unifi"
)

type wlitem struct {
	Alias    string         `json:"alias"`
	Hostname string         `json:"hostname"`
	Info     *unifi.Station `json:"-"`
}

type wlitems map[string]*wlitem
type stations map[string]*unifi.Station

type whitelist struct {
	List     wlitems
	stations stations
	client   *unifi.Client
	jsonfile string
}

func NewWhitelist(c *unifi.Client, jsonfile string) (wl *whitelist) {
	return &whitelist{
		List:     make(wlitems),
		stations: make(stations),
		client:   c,
		jsonfile: jsonfile,
	}
}

func (wl *whitelist) Load() (w *whitelist, err error) {
	var data []byte
	var list wlitems

	if data, err = ioutil.ReadFile(wl.jsonfile); err != nil {
		if _, er2 := os.Stat(wl.jsonfile); os.IsNotExist(er2) {
			err = wl.save()
		}
	} else {
		if err = json.Unmarshal(data, &list); err == nil {
			wl.List = list
		}
	}

	return wl, err
}

func (wl *whitelist) Update() (err error) {
	var ss []*unifi.Station
	var smap = make(stations)

	ss, err = wl.client.Stations(config.site)
	if err == nil {
		for _, s := range ss {
			smap[s.MAC.String()] = s
		}
		wl.stations = smap
		for mac := range wl.List {
			if s, ok := smap[mac]; ok {
				wl.List[mac].Info = s
			} else {
				wl.List[mac].Info = nil
			}
		}
	}

	return
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

func (wl *whitelist) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	self := wl.self(r)
	mode := r.URL.Query().Get("mode")
	alias := r.URL.Query().Get("alias")
	if alias == "" {
		alias = self.Hostname
	}

	switch mode {
	case "add":
		wl.add(alias, self).save()
		http.Redirect(w, r, config.base+"/", 303)
	case "rm":
		wl.rm(self).save()
		http.Redirect(w, r, config.base+"/", 303)
	default:
		http.NotFound(w, r)
	}

	return
}

func (wl *whitelist) ListHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		http.NotFound(w, r)
		return
	}

	doc := string(MustAsset("index.tmpl"))
	tpl := template.Must(template.New("index").Parse(doc))

	w.Header().Add("Content-Type", "text/html")

	tpl.Execute(
		w,
		struct {
			Info *unifi.Station
			List wlitems
		}{
			Info: wl.self(r),
			List: wl.List,
		},
	)
}

func (wl *whitelist) add(alias string, s *unifi.Station) *whitelist {
	if s != nil && s.MAC != nil {
		wl.List[s.MAC.String()] = &wlitem{
			Alias:    alias,
			Hostname: s.Hostname,
			Info:     s,
		}
	}
	return wl
}

func (wl *whitelist) assign(s *unifi.Station) *whitelist {
	if s != nil && s.MAC != nil {
		if _, ok := wl.List[s.MAC.String()]; ok {
			wl.List[s.MAC.String()].Info = s
		}
	}
	return wl
}

func (wl *whitelist) rm(s *unifi.Station) *whitelist {
	if s != nil && s.MAC != nil {
		delete(wl.List, s.MAC.String())
	}
	return wl
}

func (wl *whitelist) has(s *unifi.Station) (ok bool) {
	if s != nil && s.MAC != nil {
		_, ok = wl.List[s.MAC.String()]
	}
	return ok
}

func (wl whitelist) save() (err error) {
	var data []byte

	if data, err = json.Marshal(wl.List); err == nil {
		err = ioutil.WriteFile(wl.jsonfile, data, 0644)
	}

	if err != nil {
		fmt.Println("Error saving whitelist: " + err.Error())
	}

	return
}
