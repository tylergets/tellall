{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.tellall;
in
{
  options.services.tellall = {
    enable = mkEnableOption "TellAll service";

    # Other options your service might need.
    # port = mkOption {
    #   type = types.int;
    #   default = 8080;
    #   description = "The port on which the TellAll service runs.";
    # };
  };

  config = mkIf cfg.enable {
    # Systemd service definition
    systemd.services.tellall = {
      description = "TellAll Service";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${pkgs.tellall}/bin/tellall";
        # Replace with the actual binary path and options

        # Other service configuration
        # User = "tellall";
        # Restart = "on-failure";
      };
    };
  };
}
