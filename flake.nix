{
  description = "Filebot - program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./nix/default.nix { };
      goSDK = pkgs.go;
    in
    {
      packages.x86_64-linux.default = filebot;
      nixosModules.default = import ./nix/service.nix;
      devShells.x86_64-linux.default = pkgs.mkShell {
        hardeningDisable = [ "all" ];
        buildInputs = with pkgs; [
          # Go tools
          goSDK
          golangci-lint

          # Nix
          nixpkgs-fmt

          # Github actions 
          act
        ];

        GOROOT = "${goSDK}/share/go";
      };
    };
}
