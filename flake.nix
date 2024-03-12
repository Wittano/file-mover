{
  description = "Filebot - program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./default.nix { };
      go = pkgs.go;
    in
    {
      packages.x86_64-linux.default = filebot;
      nixosModules."filebot" = import ./service.nix;
      devShells.x86_64-linux.default = pkgs.mkShell {
        hardeningDisable = [ "all" ];
        buildInputs = with pkgs; [ go golangci-lint ];

        GOROOT = "${go}/share/go";
      };
    };
}
