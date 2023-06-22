{
  description =
    "Program to automaticlly moving files from one place to another";

  outputs = { self, nixpkgs, }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      file-mover = pkgs.callPackage ./nix { };
    in {
      defaultPackage.x86_64-linux = file-mover;
      nixosModule.default = { config }:
        import ./nix/service.nix {
          inherit config pkgs file-mover;
          lib = pkgs.lib;
        };
    };
}
