{
  description =
    "Program to automaticlly moving files from one place to another";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./default.nix { };
    in {
      packages.x86_64-linux.default = filebot;
      nixosModules."filebot" = import ./service.nix;
    };
}
