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
}
