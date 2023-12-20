{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;
  version = "v1.0.0";

  vendorHash = "sha256-yGB8VC2gcsbfcASL6pPMEVeYNtWlpMQA981BGkevL/Q=";

  patches = [
    "./patches/fix(config)__changed_path_for_nix_build.patch"
  ];
}
