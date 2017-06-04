package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mdlayher/unifi"
)

type configuration struct {
	api   string
	user  string
	pass  string
	site  string
	bind  string
	base  string
	wlist string
	freq  time.Duration
}

var (
	identifier = "mielke/1.0"
	config     = configuration{}
)

func main() {
	var c *unifi.Client
	var w *whitelist
	var err error

	fmt.Println(identifier + " initialising...")

	if c, err = initClient(); err == nil {
		if _, err = initSite(c); err == nil {
			if w, err = NewWhitelist(c, config.wlist).Load(); err == nil {
				go w.UpdateLoop(config.freq)
				fmt.Println(identifier + " now serving...")
				http.ListenAndServe(config.bind, w.Router())
				os.Exit(0)
			}
		}
	}

	fmt.Println("error: " + err.Error())
	os.Exit(1)
}

// initClient connects and logs into the unifi endpoint
func initClient() (c *unifi.Client, err error) {
	c, err = unifi.NewClient(
		config.api,
		unifi.InsecureHTTPClient(4*time.Second),
	)

	if err == nil {
		c.UserAgent = identifier
		err = c.Login(config.user, config.pass)
	}

	return
}

// initSite checks, whether the configured site name exists
func initSite(c *unifi.Client) (site *unifi.Site, err error) {
	var sites []*unifi.Site
	var siteNames []string

	if sites, err = c.Sites(); err == nil {
		for _, s := range sites {
			siteNames = append(siteNames, s.Name)
			if s.Name == config.site {
				site = s
				config.site = site.Name
				break
			}
		}

		if site == nil {
			err = fmt.Errorf(
				"could not find site '%s' in ['%s'] on '%s'",
				config.site,
				strings.Join(siteNames, "','"),
				config.api,
			)
		}
	}

	return
}

func init() {
	mkFlag(&config.api, "api", "", "UniFi API")
	mkFlag(&config.user, "user", "", "UniFi API Username")
	mkFlag(&config.pass, "pass", "", "UniFi API Password")
	mkFlag(&config.site, "site", "default", "UniFi Site")
	mkFlag(&config.bind, "bind", ":5520", "HTTP Host:Port to bind to")
	mkFlag(&config.base, "base", "/", "Reverse Proxy Path Prefix")
	mkFlag(&config.wlist, "whitelist", "whitelist.json", "List File")
	mkFlag(&config.freq, "freq", "30s", "Update Frequency")

	flag.Parse()
}

func mkFlag(v interface{}, name string, def string, desc string) {
	val := os.Getenv("MIELKE_" + strings.ToUpper(name))
	if val == "" {
		val = def
	}
	switch v.(type) {
	case *string:
		flag.StringVar(v.(*string), name, val, desc)
	case *time.Duration:
		dur, _ := time.ParseDuration(val)
		flag.DurationVar(v.(*time.Duration), name, dur, desc)
	}
}
