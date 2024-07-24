let
  # Rolling updates, not deterministic.
  pkgs = import (fetchTarball("channel:nixpkgs-unstable")) {};
in pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.air
    pkgs.gopls
    pkgs.awscli2
  ];
}
