{
  inputs = {
    nixpkgs = {
      url = "github:nixos/nixpkgs/nixpkgs-unstable";
    };
  };

  outputs = { self, nixpkgs, ... }@inputs:
  let
    pkgs = import nixpkgs { system = "x86_64-linux"; config.allowUnfree = true; };
  in
  {
    devShells.x86_64-linux.default = pkgs.mkShell {
          shellHook = "export CGO_ENABLED=1";
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.delve
            pkgs.gotools
            pkgs.golangci-lint
            pkgs.nodejs_20
            pkgs.sqlite
            pkgs.libgcc
          ];
        };
  };
}