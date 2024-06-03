{
  description = "cafédodo";

  inputs.poetry2nix.url = "github:nix-community/poetry2nix";

  outputs = { self, nixpkgs, poetry2nix, ... }:
    let
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems = [ "x86_64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {

      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          cafedodo = pkgs.buildGoModule {
            pname = "cafedodo";
            inherit version;
            src = ./service/server;

            # This hash locks the dependencies of this package. It is
            # necessary because of how Go requires network access to resolve
            # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
            # details. Normally one can build with a fake hash and rely on native Go
            # mechanisms to tell you what the hash should be or determine what
            # it should be "out-of-band" with other tooling (eg. gomod2nix).
            # To begin with it is recommended to set this, but one must
            # remember to bump this hash when your dependencies change.
            vendorHash = pkgs.lib.fakeHash;
            # vendorHash = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
          };
        });

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          inherit (poetry2nix.lib.mkPoetry2Nix { pkgs = nixpkgsFor.${system}; }) mkPoetryEnv;
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              gcc
              zlib
              go
              gopls
              gotools
              go-tools
              nodejs_20
              nodePackages.pyright
              (mkPoetryEnv { projectDir = ./checker; })
              poetry
              ruff
              enochecker-test
              vscode-langservers-extracted
              tailwindcss-language-server
              (writeScriptBin "ect" ''
                ENOCHECKER_TEST_CHECKER_ADDRESS="127.0.0.1" ENOCHECKER_TEST_CHECKER_PORT="16969" ENOCHECKER_TEST_SERVICE_ADDRESS=$(ifconfig wlp3s0 | grep 'inet ' | xargs | cut -d ' ' -f2) enochecker_test $@
              '')
            ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.cafedodo);
    };
}
