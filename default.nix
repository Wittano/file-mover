{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;

  vendorHash = "sha256-plRphEIwtPoej+bM4fChhOjBGO/BJ2KoCjZnyD/Z634=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];

  goMod = ./.;
}
