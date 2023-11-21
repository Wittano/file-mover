{ buildGoModule }:

buildGoModule rec {
  name = "filebot";
  src = ./.;

  # TODO Update vendorHash during push a relase version v1.0
  vendorHash = "sha256-XPRA1i8guYzLDEu5/QCzPhO/CHjNtgx2WNyafuoKjzc=";

  goMod = ./.;
}
