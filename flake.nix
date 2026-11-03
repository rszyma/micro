{
  description = "micro editor with pixel-precision mouse clicks";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f system);
    in
    {
      packages = forEachSystem (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = self.packages.${system}.micro;
          micro = pkgs.buildGoModule {
            pname = "micro";
            version = "2.0.15+fix-cursor-precision-on-mouse-click";

            src = pkgs.lib.cleanSource ./.;

            # vendored from go.mod/go.sum (tcell replace points at the fork)
            vendorHash = "sha256-leX7dVXe4O0QLkUaJ9/l6zMJVHw8bB3zUt3kZb5FWYI=";

            proxyVendor = true;
            subPackages = [ "cmd/micro" ];

            nativeBuildInputs = [ pkgs.installShellFiles ];

            outputs = [
              "out"
              "man"
            ];

            ldflags = [
              "-s"
              "-w"
              "-X github.com/micro-editor/micro/v2/internal/util.Version=2.0.15+fix-cursor-precision-on-mouse-click"
            ];

            preBuild = ''
              GOOS= GOARCH= go generate ./runtime
            '';

            postInstall = ''
              installManPage assets/packaging/micro.1
              install -Dm444 assets/packaging/micro.desktop $out/share/applications/micro.desktop
              install -Dm644 assets/micro-logo-mark.svg $out/share/icons/hicolor/scalable/apps/micro.svg
            '';

            meta = {
              description = "micro editor with pixel-precision mouse clicks";
              mainProgram = "micro";
            };
          };
        }
      );
    };
}
