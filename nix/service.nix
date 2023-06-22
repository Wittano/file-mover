{ config, lib, pkgs, ... }:

let
    cfg = config.services.file-mover;
    program = pkgs.callPackage ./default.nix {};
in {
  options = {
    services.file-mover = {
      enable = lib.mkEnableOption "Enable file-mover service";
      user = lib.mkOption {
        type = lib.types.str;
        example = "wittano";
        description = ''
          Specific user, who will be run service
        '';
      };
      configPath = lib.mkOption {
        type = lib.types.str;
        default = "$HOME/.config/file-mover/config.toml";
        example = "/home/wittano/config.toml";
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
    systemd.services.file-mover = {
      description = ''
        Service to automaticlly sorting files
      '';
      serviceConfig.User = "${cfg.user}";
      wantedBy = [ "multi-user.target" ];
      script = ''
        ${program}/bin/file-mover -c ${cfg.configPath} -u ${cfg.updateInterval}
      '';
    };
  };
}
