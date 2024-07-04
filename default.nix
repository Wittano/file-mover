{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;
  version = "v1.0.1";

  vendorHash = "sha256-yto8EVoOjBxmIEcxQMAcsogTlXx0A/UIAoUO+9pljoA=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];
}
