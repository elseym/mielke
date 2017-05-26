package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"strings"

	"github.com/mdlayher/unifi"
)

var (
	identifier = "mielke/1.0"

	config = struct {
		api    *string
		user   *string
		pass   *string
		site   *string
		bind   *string
		prefix *string
		wlist  *string
	}{
		flag.String("api", "", "UniFi API"),
		flag.String("user", "", "UniFi API Username"),
		flag.String("pass", "", "UniFi API Password"),
		flag.String("site", "default", "UniFi Site"),
		flag.String("bind", ":5520", "HTTP Host:Port to bind to"),
		flag.String("prefix", "", "Reverse Proxy Path Prefix"),
		flag.String("whitelist", "whitelist.json", "List File"),
	}
)

func main() {
	var c *unifi.Client
	var w *whitelist
	var err error

	flag.Parse()

	fmt.Println(identifier + " initialising...")

	if c, err = initClient(); err == nil {
		if _, err = initSite(c); err == nil {
			if w, err = NewWhitelist(c, *config.wlist).Load(); err == nil {
				go update(w)
				http.HandleFunc("/whitelist", w.UpdateHandler)
				http.HandleFunc("/", w.ListHandler)
				fmt.Println(identifier + " now serving...")
				http.ListenAndServe(*config.bind, nil)
				return
			}
		}
	}

	fmt.Println("error: " + err.Error())
	os.Exit(1)
}

func update(wl *whitelist) {
	for {
		fmt.Println("updating...")
		wl.Update()
		time.Sleep(23 * time.Second)
	}
}

// initClient connects and logs into the unifi endpoint
func initClient() (c *unifi.Client, err error) {
	c, err = unifi.NewClient(
		*config.api,
		unifi.InsecureHTTPClient(4*time.Second),
	)

	if err == nil {
		c.UserAgent = identifier
		err = c.Login(*config.user, *config.pass)
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
			if s.Name == *config.site {
				site = s
				*config.site = site.Name
				break
			}
		}

		if site == nil {
			err = fmt.Errorf(
				"could not find site '%s' in ['%s'] on '%s'",
				*config.site,
				strings.Join(siteNames, "','"),
				*config.api,
			)
		}
	}

	return
}
