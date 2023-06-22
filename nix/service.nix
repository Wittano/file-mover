{ config, lib, pkgs, file-mover, ... }:

let cfg = config.services.file-mover;
in with lib; {
  options = {
    services.file-mover = {
      enable = mkEnableOption "Enable file-mover service";
      user = mkOption {
        type = types.str;
        example = "wittano";
        description = ''
          Specific user, who will be run service
        '';
      };
      configPath = mkOption {
        type = types.str;
        default = "$HOME/.config/file-mover/config.toml";
        example = "/home/wittano/config.toml";
        description = ''
          Path to program configuration.

          REMEMBER!
          Program accept only toml-formatted files
        '';
      };
      updateInterval = kOption {
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
    systemd.services.file-mover = {
      description = ''
        Service to automaticlly sorting files
      '';
      serviceConfig.User = "${cfg.user}";
      wantedBy = [ "multi-user.target" ];
      script = ''
        ${file-mover}/bin/file-mover -c ${cfg.configPath} -u ${cfg.updateInterval}
      '';
    };
  };
}
