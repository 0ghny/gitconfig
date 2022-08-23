# gitconfigs

> Cli to manage multiple gitconfigs with ease


GitConfigs cli, helps with the tedious management of multiple gitconfigs when you need to have different configurations per location. It helps to manage the includeIf sections

```ini
[includeIf "gitdir:/a/location/"]
    path = /path/to/custom/gitconfig
```

It also offers a git config wrapper command so you can get or set configuration properties directly to specific location configurations.


## Install

Go to release page, download the binary of your platform / arch and start using it

### Install using asdf

There's a plugin for asdf to make it easy to use gitconfigs.

```shell
# Having asdf installed and configured in your system
asdf plugin add gitconfigs
asdf install gitconfigs latest
asdf global gitconfig latest
```

## Usage

```shell
Manage multiple location based git configurations easily

Usage:
  gitconfigs [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Executes git config [config-key] [config value] on an specified location
  help        Help about any command
  location    Manage a gitconfigs location
  locations   List configured locations

Flags:
  -c, --git-config string   Git configuration file (default "/home/user/.gitconfig")
  -h, --help                help for gitconfigs
  -v, --verbosity int       Verbosity level from 0 to 4
      --version             version for gitconfigs
```

### Retrieve current configured locations

```shell
# It prints all configured locations, in this example, theres only one location
# with key "github" that i use to my github specific configurations like
# changing user name and email.
$ gitconfigs locations

+---+--------+--------------------------+------------------------------------------+
| # | KEY    | LOCATION                 | GITCONFIG                                |
+---+--------+--------------------------+------------------------------------------+
| 0 | github | /code/                   | /home/user/.gitconfigs/github.gitconfig  |
+---+--------+--------------------------+------------------------------------------+
```

### Create a new location config

```shell
# This creates a new location where you are in git repositories under /code
# with key "github" with the configuration in /home/user/.gitconfigs/github.gitconfig

# We can create it in 2 ways, BEING in the directory /code, or in any other
# but passing the location pamarater

# Option 1:
$ cd /code
$ gitconfigs location new --key github

# Options 2:
$ gitconfigs location new --key github --location /code

# If you check your ~/.gitconfig file, it should shown at the very end
$ cat ~/.gitconfig
# ...
# gitconfigs.location.key github
[includeIf "gitdir:/code/"]
    path = /home/user/.gitconfigs/github.gitconfig
# ...
```

### Get or Set a configuration property on a location

Like using traditional `git config` command, this is wrapper to allow specify the location configuration file to operate.
Internally, it just wraps git config but setting environment variable `GIT_CONFIG` to the configured key location.

**Note**: if not --key parameter is provided, it will not pass GIT_CONFIG so will be a regular git config command. It is useful to check git applied configurations as they will be applied when running git commands in your system git commits, pulls, push...


```shell
# Get the user.name property in our recently created location
# It is getted from a template
$ gitconfigs config --key github user.name
anonymmous

# Sets a different user name
$ gitconfigs config --key github user.name 0ghny
# Get the new name
$ gitconfigs config --key github user.name
0ghny

# Being in /code/repository directory, we can test the git config user.name with gitconfigs as
$ cd /code/repository
$ gitconfigs config user.name
0ghny
# 0ghny will be used as your name when operate with git

```