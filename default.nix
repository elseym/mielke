with import <nixpkgs> {};
stdenv.mkDerivation ({
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
    (cd web-client && yarn install && yarn build)
    go get -u -v github.com/jteeuwen/go-bindata/...
    go get -u -v github.com/mdlayher/unifi
    ./vendor/bin/go-bindata -debug -o assets.go -prefix web-client/public web-client/public/mielke.html
    go run *.go
  '';
} // ( import ./secrets.nix ))
