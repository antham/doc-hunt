# Doc-hunt [![codecov](https://codecov.io/gh/antham/doc-hunt/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/doc-hunt) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/doc-hunt)](https://goreportcard.com/report/github.com/antham/doc-hunt) [![GitHub tag](https://img.shields.io/github/tag/antham/doc-hunt.svg)]()

Doc-hunt track changes occuring in a source code to help keeping documentation up to date.

[![asciicast](https://asciinema.org/a/6oaa4q1tfyh3x9i38m0bmkqoy.png)](https://asciinema.org/a/6oaa4q1tfyh3x9i38m0bmkqoy)

## How it works ?

### Example : README.md

This repository contains a `README.md` and it rely heavily on features made in `cmd` module (cmd folder in this repository), so to keep file `README.md` up to date, we need to track some files in`cmd` folder. You can check this repository to have an example how it works with travis.

#### Part 1 : Setup

To setup this tracking we do :

```
doc-hunt add config README.md 'cmd/(config.*|update|version|check)(?<!_test).go'
```

It creates a file `.doc-hunt` at the root of this repository and we record it to version control system, this file will track every changes occuring in files recorded in config.

#### Part 2 : Tracking changes

After few changes in source code, we want to check if we have to update `README.md`, we run :

```
doc-hunt check
```

And we get :

```
----
README.md

  Updated

    => cmd/check.go
----
```

#### Part 3 : Update doc-hunt

We check if we made changes in file `cmd/check.go` that could be documented in `README.md`. When everything looks fine, we update `.doc-hunt` file doing (we are going to commit `.doc-hunt` afterwards) :

```
doc-hunt update
```

## Install

From [release page](https://github.com/antham/doc-hunt/releases) download the binary according to your system architecture.

You can pull a docker container from dockerhub, run this command from you root path application :

`docker run -it --rm -v $PWD:/app/ antham/doc-hunt config list`

if you need a different binary version add it as tag like so :

`docker run -it --rm -v $PWD:/app/ antham/doc-hunt:2.1.1 config list`

### Warning

`doc-hunt` keep a track of version app used to record configuration, it's necessary to use a binary with same major version or configuration needs to be redo to update to a new major version.

## Usage

```bash
Ensure your documentation is up to date

Usage:
  doc-hunt [command]

Available Commands:
  check       Check if documentation update could be needed
  config      List, add or delete configuration
  update      Update documentation references
  version     App version

Flags:
  -h, --help   help for doc-hunt

Use "doc-hunt [command] --help" for more information about a command.
```

### config add

A config line is made with a document (a file, a folder in current repository or an external URL) with one or several sources (a pcre regexp matching file sources that need to be tracked).

#### Example 1

Add a new configuration on all php files in two different folders :

```
doc-hunt config add README.md folder1/.*.php,folder2/.*.php
```

#### Example 2

Add a new configuration on all php files in a folder and exclude test files :

```
doc-hunt config add README.md 'folder1/.*(?<!_test).php'
```

#### -n flag, try config

It's possible to simulate what this command will do adding a `-n` option like so :

```
doc-hunt config add -n README.md folder1/.*.php
```

### config list

List all existing configurations :

```
doc-hunt config list
```

### config del

Remove unwanted configuration (This command will launch a prompt) :

```
doc-hunt config del
```

### check

To check all changes occuring in a source code :

```
doc-hunt check
```

#### -e flag, crash when a change occurs

To make `doc-hunt` crashes when changes are detected (useful in a CI), add a `e` flag like so :

```
doc-hunt check -e
```

### update

To record every changes and make `doc-hunt` not complaining anymore :

```
doc-hunt update
```

New version of `.doc-hunt` must be committed with all changes added.
