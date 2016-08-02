Doc-hunt [![Build Status](https://travis-ci.org/antham/doc-hunt.svg?branch=master)](https://travis-ci.org/antham/doc-hunt) [![codecov](https://codecov.io/gh/antham/doc-hunt/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/doc-hunt) [![codebeat badge](https://codebeat.co/badges/dc8062aa-0b73-4d58-8b6e-a3b336409ba8)](https://codebeat.co/projects/github-com-antham-doc-hunt)
========

Doc-hunt track changes occuring in your source code to help you keeping your documentation up to date.

## How it works ?

### Example : README.md

This repository contains a ```README.md``` and it rely heavily on features made in ```cmd``` module (cmd folder in this repository), so to keep file ```README.md``` up to date, we need to track ```cmd``` folder. You can check this repository to have an example how it works with travis.

#### Part 1 : Setup

To setup this tracking we do :

```
doc-hunt add config README.md cmd
```

It creates a file ```.doc-hunt``` at the root of this repository and we record it to version control system, this file will track every changes occuring in files recorded in config.

#### Part 2 : Tracking changes

After few changes in source code, we want to check if we have to update ```README.md```, we run :

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

We check if we made changes in file ```cmd/check.go``` that could be documented in ```README.md```. When everything looks fine, we update ```.doc-hunt``` file doing (we are going to commit ```.doc-hunt``` afterwards) :

```
doc-hunt update
```

## Install

From [release page](https://github.com/antham/doc-hunt/releases) download the binary according to your system architecture

### Warning

```doc-hunt``` keep a track of version app you used to record you configuration, you need to keep a binary with same major version or you have to redo you configuration if you want to update to a new major version.

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

### Config

A config line is made with a document (a file, a folder in current repository or an external URL) with one or several sources (a regexp matching file sources you want to track).

To add a new configuration on all php files in two different folders run :

```
doc-hunt config add README.md folder1/.*.php,folder2/.*.php
```

You can record as many configurations you want and list them with :

```
doc-hunt config list
```

You can remove unwanted configuration through (This command will launch a prompt) :

```
doc-hunt config del
```

### Check

To check all changes occuring on your source code run :

```
doc-hunt check
```

If you want to make it crashes everytime changes are detected, it's useful in a CI, just add -e flags like so :

```
doc-hunt check -e
```

### Update

To record every changes run :

```
doc-hunt update
```

Commit this new version of .doc-hunt in your repository.
