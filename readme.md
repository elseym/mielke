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

* use the [.env.dist](env.dist) as template to create a `.env` file that contains the secrets
  * `MIELKE_API` set to the unifi api. example: http://unifi.local:8443/
  * `MIELKE_USER` set to the unifi api user
  * `MIELKE_PASS` set to the unifi api password
* run nix-shell to get a shell with all dependencies `nix-shell`

or

* [build the web-client frontend](web-client/readme.md)
* install dependency `go get -u -v github.com/jteeuwen/go-bindata/...`
* install dependency `go get -u -v github.com/mdlayher/unifi`
* generate assets with `go-bindata -debug -o assets.go -prefix web-client/public web-client/public/mielke.html`;
  run without the `-debug` switch before committing assets
* start mielke with `go run *.go`

### template data
see `whitelist.go`
