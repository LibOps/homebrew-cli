# libops

Command line utility to interact with your LibOps site.

## Install

First, you must install Google Cloud's CLI [gcloud](https://cloud.google.com/sdk/docs/install)

### Homebrew
You can install the LibOps CLI using homebew
```
brew tap libops/libops
brew install libops
```

### Download Binary

Instead of homebrew, you can download a binary for your system from [the latest release](https://github.com/LibOps/cli/releases/latest)

Then put the binary in a directory that is in your `$PATH`

## Usage

After installation, the utility must be ran from within the locally checked out repository that contains your site's source code

```
$ git clone git@github.com:libops/your-site-repo
$ cd your-site-repo
$ libops drush -- status

Drupal version   : 9.5.8
Site URI         : http://default
DB driver        : mysql
DB hostname      : mariadb
DB port          :
DB username      : development
DB name          : drupal
Database         : Connected
Drupal bootstrap : Successful
Default theme    : libops_www
Admin theme      : claro
PHP binary       : /usr/local/bin/php
PHP OS           : Linux
PHP version      : 8.1.17
Drush script     : /code/vendor/bin/drush
Drush version    : 11.5.1
Drush temp       : /tmp
Drush configs    : /code/vendor/drush/drush/drush.yml
Install profile  : standard
Drupal root      : /code/web
Site path        : sites/default
Files, Public    : sites/default/files
Files, Private   : /private
Files, Temp      : /tmp
```
