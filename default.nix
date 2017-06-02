with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "mielke";
  src = ./.;
  buildInputs = [
    git
    nodejs
    go
    yarn
  ];
  GOPATH = "$(pwd)/vendor";
  shellHook = ''
    [ -f .env ] && . .env || . .env.dist
    (cd web-client && yarn install && yarn build:dev &)
    go get -u -v github.com/jteeuwen/go-bindata/...
    go get -u -v github.com/mdlayher/unifi
    ./vendor/bin/go-bindata -debug -o assets.go -prefix web-client/public web-client/public/mielke.html
    go run *.go
  '';
}
