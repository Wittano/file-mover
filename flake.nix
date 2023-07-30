{
  description =
    "Program to automaticlly moving files from one place to another";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./nix { };
    in {
      defaultPackage.x86_64-linux = filebot;
      nixosModules."filebot" = import ./nix/service.nix;
    };
}
