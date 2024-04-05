{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_22
    pkgs.git
    pkgs.gnumake
  ];

  shellHook = ''
    echo "go:       $(go version)"
    echo "git:      $(git --version)"
    echo "make:     $(make --version | head -n 1)"
    '';
}
