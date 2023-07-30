{
  description =
    "Program to automaticlly moving files from one place to another";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      file-mover = pkgs.callPackage ./nix { };
    in rec {
      defaultPackage.x86_64-linux = file-mover;
      nixosModules."filebot" = import ./nix/service.nix;
    };
}
