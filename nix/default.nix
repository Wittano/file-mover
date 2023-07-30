{ buildGoModule }:

buildGoModule rec {
  name = "file-mover";
  src = ./../.;

  vendorHash = "sha256-YVbhTJ1gwwpWhxgUHQlp+udSx3sLtMlr1TiZWsIeORA";

  goMod = ./.;
}
