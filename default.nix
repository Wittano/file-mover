{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;
  version = "v1.0.0";

  vendorHash = "sha256-plRphEIwtPoej+bM4fChhOjBGO/BJ2KoCjZnyD/Z634=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];

  goMod = ./.;
}
