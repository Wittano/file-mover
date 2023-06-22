{ buildGoModule }:

buildGoModule rec {
  name = "file-mover";
  src = ./../.;

  vendorHash = "sha256-YVbhTJ1gwwpWhxgUHQlp+udSx3sLtMlr1TiZWsIeORA";

  postInstall = ''
    mv $out/bin/src $out/bin/file-mover
  '';

  goMod = ./.;
}
