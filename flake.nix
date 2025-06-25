{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      inherit (nixpkgs) lib;
      forAllSystems = lib.genAttrs lib.systems.flakeExposed;
    in
    {
      nixosModules.networkd-ipmon = ./modules/networkd-ipmon.nix;
      nixosModules.default = self.nixosModules.networkd-ipmon;
      darwinModules.networkd-ipmon = ./modules/networkd-ipmon.nix;
      darwinModules.default = self.darwinModules.networkd-ipmon;

      packages = forAllSystems (system: {
        networkd-ipmon = nixpkgs.legacyPackages.${system}.callPackage ./pkgs/networkd-ipmon.nix { };
        default = self.packages.${system}.networkd-ipmon;
      });
      overlays.default = final: prev: {
        networkd-ipmon = prev.callPackage ./pkgs/networkd-ipmon.nix { };
      };
      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt-rfc-style);

      devShells = forAllSystems (system: {
        default =
          let
            pkgs = nixpkgs.legacyPackages.${system};
            packages = with pkgs; [
              go_1_24
            ];
          in
          derivation {
            name = "shell";
            inherit system packages;
            builder = "${pkgs.bash}/bin/bash";
            outputs = [ "out" ];
            stdenv = pkgs.writeTextDir "setup" ''
              set -e

              for p in $packages; do
                PATH=$p/bin:$PATH
              done
            '';
          };
      });
    };
}
