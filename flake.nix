{
  description = "A Nix environoment for Lizzy Wizzy discord Bot";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }: let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [ go ffmpeg ];

      shellHook = ''
        export GOPATH=$HOME/go/bin:$PATH
        export $(grep -v '^#' .env | xargs)
      '';
    };
  };
}
