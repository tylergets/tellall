flake: { options, config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.tellall;

  tellall = flake.packages.${pkgs.stdenv.hostPlatform.system}.default;
  hmModule = flake.homeManagerModule;
in
{
  config = {
    home-manager.sharedModules = [
      flake.homeManagerModules.tellall
    ];
  };
}
