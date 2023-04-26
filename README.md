# LibOps CLI

Command line utility to interact with your LibOps site.

## Install

First, you must install Google Cloud's CLI [gcloud](https://cloud.google.com/sdk/docs/install)

After `gcloud` is installed you can install LibOps CLI using homebrew or a binary.

### Homebrew
You can install the LibOps CLI using homebew
```
brew tap libops/cli
brew install libops
```

### Download Binary

Instead of homebrew, you can download a binary for your system from [the latest release](https://github.com/LibOps/homebrew-cli/releases/latest)

Then put the binary in a directory that is in your `$PATH`


## Usage

After installation, LibOps CLI must be ran from within the locally checked out repository that contains your site's source code

```
$ git clone git@github.com:libops/your-site-repo
$ cd your-site-repo
$ libops --help                       
Interact with your libops site

Usage:
  libops [command]

Available Commands:
  backup      Backup your libops environment
  completion  Generate the autocompletion script for the specified shell
  config-ssh  Populate ~/.ssh/config with LibOps development environment
  drush       Run drush commands on your libops environment
  help        Help about any command
  sync-db     Transfer the database from one environment to another

Flags:
  -e, --environment string   LibOps environment (default "development")
  -h, --help                 help for libops
  -p, --site string          LibOps project/site (default "your-site-repo")
  -v, --version              version for libops

Use "libops [command] --help" for more information about a command.
```

## Updating

Updating your LibOps CLI depends on the method used to install.

### Homebrew

If homebrew was used, you can simply upgrade the homebrew formulae for LibOps CLI

```
brew update && brew upgrade libops
```

### Download Binary

If the binary was downloaded and added to the `$PATH` updating could look as follows. Requires `curl`, `tar`, `jq`

```
# update for your architecture
ARCH="homebrew-cli_Linux_x86_64.tar.gz"
curl -s https://api.github.com/repos/LibOps/homebrew-cli/releases/latest > latest.json
URL=$(jq -rc '.assets[] | select(.name == "'$ARCH'") | .browser_download_url' latest.json)
echo "Fetching latest libops CLI release from $URL"
curl -Ls -o $ARCH "$URL"
tar -zxvf $ARCH
mv libops /usr/local/bin/
rm $ARCH LICENSE README.md
```
