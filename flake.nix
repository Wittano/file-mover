{
  description =
    "Filebot - program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in rec {
      packages.x86_64-linux.default = pkgs.buildGoModule {
        name = "filebot";
        src = ./.;
        version = "v1.0.0";

        vendorHash = "sha256-zga1pCBqisDLzDN6rO68iCQlGXmTfkUk+fqNI54yhNo=";

        patches = [ "./patches/fix(config)__changed_path_for_nix_build.patch" ];
      };
      nixosModules."filebot" = pkgs.callPackage ./service.nix { filebot = packages.x86_64-linux.default; };
      devShells.x86_64-linux.default = pkgs.mkShell {
        hardeningDisable = [ "all" ];
        buildInputs = with pkgs; [
          # Go tools
          go
          golangci-lint

          # Nix
          nixfmt

          # Github actions 
          act
        ];

        GOROOT = "${pkgs.go}/share/go";
      };
    };
}
