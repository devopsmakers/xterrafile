# XTerrafile [![Build Status](https://circleci.com/gh/devopsmakers/xterrafile.svg?style=shield)](https://circleci.com/gh/devopsmakers/xterrafile)[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdevopsmakers%2Fxterrafile.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdevopsmakers%2Fxterrafile?ref=badge_shield)

`xterrafile` is a binary written in Go to manage external modules from Github for use with (but not limited to) Terraform. See this [article](http://bensnape.com/2016/01/14/terraform-design-patterns-the-terrafile/) for more information on how it was introduced in a Ruby rake task.

Inspired by:
* https://github.com/coretech/terrafile
* https://github.com/claranet/python-terrafile

## How to install

### Requirements
* `git`

### macOS

```sh
brew tap devopsmakers/xterrafile && brew install xterrafile
```

### Linux
Download your preferred flavour from the [releases](https://github.com/devopsmakers/xterrafile/releases/latest) page and install manually.

For example:
```sh
curl -L https://github.com/devopsmakers/xterrafile/releases/download/v{VERSION}/xterrafile_{VERSION}_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin
```

## How to use
By default, `xterrafile` expects a file named `./Terrafile` which will contain your terraform module dependencies in YAML format.

Specifying modules in your `Terrafile`:
```
# Terraform Registry module
terraform-digitalocean-droplet:
  source: "terraform-digitalocean-modules/droplet/digitalocean"
  version: "0.1.7" // If version is empty, the latest will be fetched

# Git module (HTTPS)
terraform-digitalocean-droplet:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  // No version will checkout default branch (usually master)

# Git module (SSH + Tag)
terraform-digitalocean-droplet:
  source: "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "v0.1.7" // Checkout tags

# Git module (HTTPS + Branch as url parameter)
terraform-digitalocean-droplet:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git?ref=new_feature"

# Git module (SSH + Commit)
terraform-digitalocean-droplet:
  source: "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "2e6b9729f3f6ea3ef5190bac0b0e1544a01fd80f" // Checkout a commit

# Get a path from within a Git monorepo
terraform-digitalocean-droplet:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"
  version: "v0.1.7"
  path: "examples/simple"

# Get a path from within a Git monorepo - alternate syntax
terraform-digitalocean-droplet:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git//examples/simple?ref=v0.1.7"

# Compressed archive (extracting a directory from inside archive)
terraform-digitalocean-droplet:
  source: "https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet/archive/v0.1.7.tar.gz//terraform-digitalocean-droplet-0.1.7"

# Local directory module
terraform-digitalocean-droplet:
  source: "../../modules/terraform-digitalocean-droplet"
```

You can specify modules using Terraform's `source` specifications:
https://www.terraform.io/docs/modules/sources.html

The `version` can be a tag, a branch or a commit hash. By default, `xterrafile`
will checkout the default branch of a module which is usually `master`.

Modules will be downloaded to `./vendor/xterrafile/` by default.

### Example Usage
Help:
```
xterrafile help
Manage vendored modules with a YAML file.

Usage:
  xterrafile [command]

Available Commands:
  help        Help about any command
  install     Installs the modules in your Terrafile
  version     Show version information for xterrafile

Flags:
  -d, --directory string   module directory (default "vendor/xterrafile")
  -f, --file string        config file (default "Terrafile")
  -h, --help               help for xterrafile

Use "xterrafile [command] --help" for more information about a command.
```

Defaults:
```
xterrafile install
```

Custom "Terrafile":
```
xterrafile -f Saltfile install
```

Custom "Terrafile" and custom download directory:
```
xterrafile -f Saltfile -d /srv/formulas install
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fdevopsmakers%2Fxterrafile.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fdevopsmakers%2Fxterrafile?ref=badge_large)
