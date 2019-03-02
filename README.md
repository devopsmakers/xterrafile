# Terrafile [![Build Status](https://circleci.com/gh/devopsmakers/xterrafile.svg?style=shield)](https://circleci.com/gh/devopsmakers/xterrafile)

`xterrafile` is a binary written in Go to manage external modules from Github for use in (but not limited to) Terraform. See this [article](http://bensnape.com/2016/01/14/terraform-design-patterns-the-terrafile/) for more information on how it was introduced in a Ruby rake task.

## How to install

### macOS

```sh
brew tap devopsmakers/xterrafile && brew install xterrafile
```

### Linux
Download your preferred flavor from the [releases](https://github.com/devopsmakers/xterrafile/releases/latest) page and install manually.

For example:
```sh
curl -L https://github.com/devopsmakers/xterrafile/releases/download/v{VERSION}/terrafile_{VERSION}_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin
```

## How to use
By default, `xterrafile` expects a file named `Terrafile` which will contain your terraform module dependencies in YAML format.

An example Terrafile:
```
tf-aws-vpc:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "v1.46.0"
tf-aws-vpc-experimental:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "master"
tf-aws-vpc-default:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
tf-aws-vpc-commit:
    source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc"
    version: "01601169c00c68f37d5df8a80cc17c88f02c04d0"
```
The `version` can be a tag, a branch or a commit hash. By default, `xterrafile`
will checkout the master branch of a module.

Modules will be downloaded to `./vendor/xterrafile/`.

### Example Usage
Defaults:
```
xterrafile install
```

Custom "Terrafile":
```
xterrafile -f Saltfile
```

Custom download directory:
```
xterrafile -f Saltfile -d /srv/formulas
```
