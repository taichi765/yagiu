{
  description = "Small CLI to manage local git repositories";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs allSystems (
          system:
          f {
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      packages = forAllSystems (
        { pkgs }: {
          default = pkgs.buildGoModule {
            pname = "yagiu";
            version = "0.0.0";
            src = ./.;
            vendorHash = "sha256-q3GqoDNsqvb++E1Js5oK3+t62kIT9OxLJFu7bjyfCFY=";
          };
        }
      );
    };
}
