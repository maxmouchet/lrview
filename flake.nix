{
  description = "Zero-configuration Lightroom catalog viewer";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "nixpkgs/nixos-22.05";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        packages = {
          lrview = pkgs.buildGo118Module {
            pname = "lrview";
            version = "0.1.0";
            src = ./.;
            vendorSha256 = "sha256-9mrq1BVVDoJ+nPQkvbQM2LENUPaxGUCSCrJJ2Dnb550=";
          };
        };
        defaultPackage = self.packages.${system}.lrview;
      }
    );
}
