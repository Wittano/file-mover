{
  description = "Filebot - program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in
    {
      packages.x86_64-linux.default = pkgs.callPackage ./default.nix { };
      nixosModules.default = ./service.nix;
      devShells.x86_64-linux.default = pkgs.callPackage ./shell.nix { };
    };
}
