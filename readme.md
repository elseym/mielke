# mielke

mielke queries your unifi api and presents whitelisted clients'
online state. clients whitelist themselves by opting-in to have
their device's mac address stored and their online state shown.

a use case could be a website presenting a list of (whitelisted)
colleagues currently present in the office.

## installation

* install via `go get -u -v github.com/elseym/mielke`
* test with `$GOPATH/bin/mielke --help`

## running mielke
 * execute `mielke --help` to see which options are available
 * supply at least `api`, `user`, and `pass`
 * make sure that `whitelist` is read/writeable;
   the file will be created if it does not exist
 * access the webinterface via your nic's ip

## contributing

* install dependency `go get -u -v github.com/jteeuwen/go-bindata/...`
* install dependency `go get -u -v github.com/mdlayher/unifi`
* generate assets with `go-bindata -debug -o assets.go -prefix assets/ assets/`;
  run without the `-debug` switch before committing assets
* start mielke with `go run *.go`

### template data
```
  . (dot)
  `- List (the whitelisted clients, key: mac address) map[string]struct
    `- Alias (the user supplied name) string
    `- Hostname (the hostname at the time of whitelisting) string
    `- Online (whether this client is online) bool
    `- Info (the whitelisted station's info) *mdlayher/unifi.Statiom
  `- Info (the requesting stations's info) *mdlayher/unifi.Statiom
```
