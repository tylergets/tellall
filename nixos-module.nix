flake: { config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.tellall;

  tellall = flake.packages.${pkgs.stdenv.hostPlatform.system}.default;

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

    package = mkOption {
      type = types.package;
      default = tellall;
      description = ''
        The Tellall package to use with the service.
      '';
    };
  };

  config = mkIf cfg.enable {
    # Systemd service definition
    systemd.services.tellall = {
      description = "TellAll Service";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${lib.getBin tellall}/bin/tellall";
        # Replace with the actual binary path and options

        # Other service configuration
         User = "tellall";
         Restart = "on-failure";
         StateDirectory = "tellall";
      };
    };
  };
}
