{
  description = "bump rockets, make whoosh";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    templ.url = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, templ, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          go-capnp = pkgs.fetchFromGitHub {
            name = "go-capnp";
            owner = "capnproto";
            repo = "go-capnp";
            rev = "main";
            hash = "sha256-P6YP5b5Bz5/rS1ulkt1tSr3mhLyxxwgCin4WRFErPGM=";
          };
        in rec
        {
          capnpc-go = pkgs.buildGoModule {
            name = "capnpc-go";
            pname = "capnpc-go";
            src = go-capnp;
            sourceRoot = "go-capnp";
            vendorHash = "sha256-DRNbv+jhIHzoBItuiQykrp/pD/46uysFbazLJ59qbqY=";
            buildPhase = ''
              go install ./capnpc-go
            '';
          };
          rocket_crash = pkgs.stdenv.mkDerivation {
            name = "rocket_crash";
            pname = "rocket_crash";
            version = builtins.substring 0 8 (self.lastModifiedDate or "19700101");
            srcs = [
              go-capnp
            ];

            GOMAXPROCS = "1";

            sourceRoot = "rocket_crash-src/";
            shellHook = ''
              echo "Entered devshell"
              export CAPNPC_GO_STD="${capnpc-go.src}/std" 
            '';
            configureFlags = [
              "--with-go-capnp=../go-capnp"
            ];

            nativeBuildInputs = with pkgs; [
              go
              gotools
              gopls
              capnproto
              capnpc-go
              hurl
              templ.packages.${system}.templ
            ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.rocket_crash);
    };
}
