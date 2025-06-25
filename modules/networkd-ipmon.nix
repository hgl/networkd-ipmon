{
  config,
  pkgs,
  lib,
  ...
}:
let
  cfg = config.services.networkd-ipmon;
  concatMapAttrsToList = f: attrs: lib.concatLists (lib.mapAttrsToList f attrs);
  rulesDir =
    rules:
    pkgs.linkFarm "networkd-ipmon-rules" (
      concatMapAttrsToList (name: rule: [
        {
          name = "${name}.json";
          path = pkgs.writeText "${name}.json" (builtins.toJSON { inherit (rule) interfaces properties; });
        }
        {
          inherit name;
          path = rule.script;
        }
      ]) rules
    );
in
{
  options = {
    services.networkd-ipmon = {
      enable = lib.mkEnableOption "networkd-ipmon";
      package = lib.mkOption {
        type = lib.types.package;
        default = pkgs.callPackage ../pkgs/networkd-ipmon.nix { };
        defaultText = lib.literalExpression "pkgs.callPackage ../pkgs/networkd-ipmon.nix { }";
        description = ''
          networkd-ipmon package to use.
        '';
      };
      rules = lib.mkOption {
        type = lib.types.attrsOf (
          lib.types.submodule {
            options = {
              interfaces = lib.mkOption {
                type = lib.types.nonEmptyListOf lib.types.nonEmptyStr;
              };
              properties = lib.mkOption {
                type = lib.types.nonEmptyListOf (
                  lib.types.enum [
                    "IPV6_ADDRS"
                    "IPV4_ADDRS"
                    "PD_ADDRS"
                  ]
                );
              };
              script = lib.mkOption {
                type = lib.types.path;
              };
            };
          }
        );
      };
    };
  };
  config = {
    systemd.services.networkd-ipmon = {
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      serviceConfig = {
        ExecStart = "${cfg.package}/bin/networkd-ipmon ${rulesDir cfg.rules}";
        Restart = "on-failure";
      };
    };
  };
}
