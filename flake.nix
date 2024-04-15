{
  description = "Filebot - program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./default.nix { };
      goSDK = pkgs.go;
    in
    {
      packages.x86_64-linux.default = filebot;
      nixosModules.default = import ./service.nix;
      devShells.x86_64-linux.default = import ./shell.nix { inherit pkgs goSDK; }; 
    };
}
