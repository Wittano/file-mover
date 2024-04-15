{ pkgs ? import <nixpkgs> { }, goSDK }: pkgs.mkShell {
  hardeningDisable = [ "all" ];
  nativeBuildInputs = with pkgs; [
    # Go tools
    goSDK

    # Nix
    nixpkgs-fmt
    nixd
  ];

  GOROOT = "${goSDK}/share/go";
}
