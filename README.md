Doc-hunt [![Build Status](https://travis-ci.org/antham/doc-hunt.svg?branch=master)](https://travis-ci.org/antham/doc-hunt) [![codecov](https://codecov.io/gh/antham/doc-hunt/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/doc-hunt) [![codebeat badge](https://codebeat.co/badges/dc8062aa-0b73-4d58-8b6e-a3b336409ba8)](https://codebeat.co/projects/github-com-antham-doc-hunt)
========

Doc-hunt track changes occuring on your source code to help you keeping your documentation up to date.

## How it works ?

### Example : README.md

This repository contains a README.md and it rely heavily on features made in cmd module, so to keep my README up to date, I need to track cmd folder.

#### Part 1 : Setup

To setup this tracking I will do :

```
doc-hunt add config README.md cmd
```

It will create a file .doc-hunt at the root of my repository and I record it to my version control system, this file will track every changes occuring on files recorded in my config.

#### Part 2 : Tracking changes

After few changes on my source code, I want to check if I have to update my doc, I run :

```
doc-hunt check
```

And I get :

```
----
README.md

  Updated

    => cmd/check.go
----
```

#### Part 3 : Update doc-hunt

I will check if I made changes in cmd/check.go that could be documented in my readme. When everything looks fine, I update my doc-hunt file doing (I commit my .doc-hunt afterwards) :

```
doc-hunt update
```

## Install

From [release page](https://github.com/antham/doc-hunt/releases) download the binary according to your system architecture

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

A config line is made with a document (a file, a folder in current repository or an external URL) with one or several sources (files or folder in current repository).

To add a new configuration on a file and two folder for instance, run :

```
doc-hunt config add README.md file1.php,folder1,folder2
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
