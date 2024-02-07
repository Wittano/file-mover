{
  description = "Program to automaticlly moving files from one place to another";

  inputs.nixpkgs.url = "nixpkgs/nixos-23.11";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      filebot = pkgs.callPackage ./default.nix { };
    in
    {
      packages.x86_64-linux.default = filebot;
      nixosModules."filebot" = import ./service.nix;
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [ go ];
      };
    };
}
