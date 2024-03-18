{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;
  version = "v1.0.0";

  vendorHash = "sha256-zga1pCBqisDLzDN6rO68iCQlGXmTfkUk+fqNI54yhNo=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];
}
