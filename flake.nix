{
  description = "A basic flake with a shell";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ go_1_18 ];
      };
      packages = {
	      default = pkgs.buildGoModule {
	        pname = "tat";
	        version = "0.1";
	        src = ./.;
	        vendorSha256 = "sha256-NknIeFrOCGzEM5gCPJ0JrFXzcKSbd5zdHw+1Nb0uhSw=";
	      };
      };
    });
}
