{ mkShell, go, nixpkgs-fmt, nixd }: mkShell {
  hardeningDisable = [ "all" ];
  nativeBuildInputs = [
    # Go tools
    go

    # Nix
    nixpkgs-fmt
    nixd
  ];

  GOROOT = "${go}/share/go";
}
