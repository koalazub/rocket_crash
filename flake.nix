{
  description = "overengineered rocket crash";

  # Nixpkgs / NixOS version to use.
  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    templ.url = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, templ, ... }:
    let
      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      # Provide some binary packages for selected system types.
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
          templPkg = templ;
        in rec
        {
          capnpc-go = pkgs.buildGoModule {
            name = "capnpc-go";
            pname = "capnpc-go";
            src = go-capnp;
            sourceRoot = "go-capnp";
            vendorSha256 = "sha256-DRNbv+jhIHzoBItuiQykrp/pD/46uysFbazLJ59qbqY=";
            buildPhase = ''
              go install ./capnpc-go
            '';
          };
          staging = pkgs.stdenv.mkDerivation {
            pname = "rocket_crash";
            version = builtins.substring 0 8 (self.lastModifiedDate or "19700101");
            srcs = [ go-capnp ]; # This should be a path to the source or a derivation

            # sourceRoot should be a string pointing to the directory to use as the root for the build.
            sourceRoot = "src"; 

            preConfigure = ''
              export XDG_CACHE_HOME=$TMPDIR/.cache
              export GOPATH=$XDG_CACHE_HOME/go
            '';

            configureFlags = [
              "--with-go-capnp=../go-capnp" # This flag should be appropriate for your build
            ];

            nativeBuildInputs = with pkgs; [
              capnproto
              capnpc-go # This should be a derivation, not a flag or a string.
              go
              gopls
              # ... other inputs
              templPkg # Make sure 'templ' is a derivation too.
            ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.staging);
    };
}
