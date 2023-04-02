
with import <nixpkgs> { };
mkShell {
    buildInputs = with pkgs; [ go gopls nixfmt rnix-lsp gnumake ];
}