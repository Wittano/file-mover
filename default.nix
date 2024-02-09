{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;
  version = "v1.0.0";

  vendorHash = "sha256-zVQx3qaYqPlxirxX/FeuGHfKcsmvgWvLrKtQLkzVkgE=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];
}
