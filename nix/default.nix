{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./../.;

  vendorHash = "sha256-YVbhTJ1gwwpWhxgUHQlp+udSx3sLtMlr1TiZWsIeORA";

  goMod = ./.;
}
