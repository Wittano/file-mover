{ config, lib, pkgs, ... }:
with lib;
let
  cfg = config.services.filebot;
  filebot = pkgs.callPackage ./default.nix { };

  buildConfig = conf:
    let
      destAttrs = assert conf ? dest && (conf ? moveToTrash && conf.moveToTrash == true);
        if conf ? dest
        then { dest = conf.dest; }
        else { moveToTrash = conf.moveToTrash; };
      afterAttrs = assert conf ? after && conf.moveToTrash != true; { after = conf.after; };
    in
    {
      name = conf.name;
      value = {
        src = conf.src;
        exceptions = mkIf (conf ? exceptions) conf.exceptions;
        recursive = mkIf (conf ? recursive) conf.recursive;
        uid = conf.uid;
        gid = conf.gid;
        isRoot = assert cfg.user == "root"; conf.isRoot;
      } // destAttrs // afterAttrs;
    };

  configs = builtins.listToAttrs (builtins.map buildConfig cfg.settings);

  filebotConfig = (pkgs.formats.toml { }).generate "config.toml" configs;
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
        type = types.path;
        default = "$HOME/.setting/filebot/setting.toml";
        example = "/home/wittano/setting.toml";
        description = ''
          Path to program configuration.

          REMEMBER!
          Program accept only toml-formatted files
        '';
      };
      settings = mkOption {
        type = with types; listOf (submodule {
          options = {
            name = mkOption {
              type = str;
              example = "File";
              description = "Single entity of configuration. It's only nesessary for generatation config. It isn't affect on filebot";
            };
            isRoot = mkEnableOption "Set ownership files to root(required root permission)";
            recursive = mkEnableOption "Observe files in directories inside src path";
            moveToTrash = mkEnableOption "Move files to Trash directory. This option required abset 'dest' option";
            after = mkOption {
              type = ints.positive;
              description = "Number of days, after files should move to Trash directory";
            };
            exceptions = mkOption {
              type = listOf str;
              example = [ "files.pdf" "files.ini" ];
              description = "List of filenames, that should be ignored by filebot";
            };
            uid = mkOption {
              type = ints.positive;
              description = "New owner UID";
            };
            gid = mkOption {
              type = ints.positive;
              description = "New owner GID";
            };
            src = mkOption {
              type = listOf str;
              example = [ "/home/user/Documents" ];
              description = "List of sources files, which should be observer";
            };
            dest = mkOption {
              type = nullOr str;
              example = "/home/user/Destination";
              description = "Path to destination directory. This option can't set with moveToTrash";
            };
          };
        });
        description = "Configurtion for filebot";
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

  assertions = [
    {
      assertion = cfg.settings != null && cfg.configPath != null;
      message = "Option services.filebot.settings and services.filebot.configPath can't be set at the same time";
    }
    {
      assertion = cfg.settings == null && cfg.configPath == null;
      message = "One of option: services.filebot.settings and services.filebot.configPath must be set";
    }
  ];

  config = mkIf cfg.enable {
    systemd.services.filebot = {
      description = "Service to automaticlly sorting files";
      serviceConfig.User = mkIf (cfg.user != "root") cfg.user;
      wantedBy = [ "multi-user.target" ];
      path = [ filebot ];
      script = "filebot -c ${filebotConfig} -u ${cfg.updateInterval}";
    };
  };
}
