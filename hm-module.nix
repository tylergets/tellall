flake: { config, lib, pkgs, ... }:

let

  cfg = config.services.tellall;
  tellall = flake.packages.${pkgs.stdenv.hostPlatform.system}.default;

in {
  meta.maintainers = with lib.maintainers; [ kranzes ];

  options = {
    services.tellall = {
      enable = lib.mkEnableOption "Tellall Daemon";

      package = tellall;

      config = lib.mkOption {
        type = with types; attrsOf (oneOf [ str int bool ]);
        default = {};
      };

      extraArgs = lib.mkOption {
        type = with lib.types; listOf str;
        default = [ ];
        description = ''
          Extra arguments to be passed to the tellall executable.
        '';
      };
    };
  };

  config = lib.mkIf cfg.enable {
    assertions = [
      (lib.hm.assertions.assertPlatform "services.tellall" pkgs
        lib.platforms.linux)
    ];

    systemd.user.services.tellall = {
      Unit = {
        Description = "Tellall - MQTT notification daemon";
        After = [ "graphical-session-pre.target" ];
        PartOf = [ "graphical-session.target" ];
      };

      Service = {
        Type = "simple";
        ExecStart =  "${lib.getExe tellall} ${lib.escapeShellArgs cfg.extraArgs}";
        Restart = "on-failure";
        RestartSec = 5;
      };

      Install.WantedBy = [ "graphical-session.target" ];
    };
  };
}