{ config, lib, pkgs, ... }:
with lib;
let
  cfg = config.services.filebot;
  filebot = pkgs.callPackage ./default.nix { };
in
{
  options = {
    services.filebot = {
      enable = mkEnableOption "Enable filebot service";
      user = mkOption {
        type = types.str;
        example = "wittano";
        description = "Specific user, who will be run service";
      };
      configPath = mkOption {
        type = with types; path;
        example = ./home/wittano/setting.toml;
        description = ''
          Path to program configuration.

          REMEMBER!
          Program accept only toml-formatted files
        '';
      };
      updateInterval = mkOption {
        type = types.str;
        default = "10m";
        example = "1h";
        description = ''
          Specific time, that program updates configuration ale observed files.
          Program accept time in go standard format(see https://pkg.go.dev/time#ParseDuration. This page include acceptable time format).
        '';
      };
    };
  };

  config = mkIf cfg.enable {
    systemd.services.filebot = {
      description = "Service to automaticlly sorting files";
      serviceConfig.User = mkIf (cfg.user != "root") cfg.user;
      wantedBy = [ "multi-user.target" ];
      path = [ filebot ];
      script = "filebot -c ${cfg.configPath} -u ${cfg.updateInterval}";
    };
  };
}
