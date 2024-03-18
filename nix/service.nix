{ config, lib, pkgs, ... }:

let
  cfg = config.services.filebot;
  program = pkgs.callPackage ./default.nix { };
in
{
  options = {
    services.filebot = {
      enable = lib.mkEnableOption "Enable filebot service";
      user = lib.mkOption {
        type = lib.types.str;
        example = "wittano";
        description = ''
          Specific user, who will be run service
        '';
      };
      configPath = lib.mkOption {
        type = lib.types.path;
        default = "$HOME/.setting/filebot/setting.toml";
        example = "/home/wittano/setting.toml";
        description = ''
          Path to program configuration.

          REMEMBER!
          Program accept only toml-formatted files
        '';
      };
      updateInterval = lib.mkOption {
        type = lib.types.str;
        default = "10m";
        example = "1h";
        description = ''
          Specific time, that program updates configuration ale observed files.
          Program accept time in go standard format(see https://pkg.go.dev/time#ParseDuration. This page include acceptable time format).
        '';
      };
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.filebot = {
      description = ''
        Service to automaticlly sorting files
      '';
      serviceConfig.User = "${cfg.user}";
      wantedBy = [ "multi-user.target" ];
      script = ''
        ${program}/bin/filebot -c ${builtins.toString cfg.configPath} -u ${cfg.updateInterval}
      '';
    };
  };
}
