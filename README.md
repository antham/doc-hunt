Doc-hunt [![Build Status](https://travis-ci.org/antham/doc-hunt.svg?branch=master)](https://travis-ci.org/antham/doc-hunt) [![codecov](https://codecov.io/gh/antham/doc-hunt/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/doc-hunt) [![codebeat badge](https://codebeat.co/badges/dc8062aa-0b73-4d58-8b6e-a3b336409ba8)](https://codebeat.co/projects/github-com-antham-doc-hunt)
========

Doc-hunt track changes occuring on your source code to help your keeping your documentation up to date.

## How it works ?

### Example : README.md

For instance, this repository contains a README.md. This readme is really tied to cmd module that live in this repository cause it's there where we have all commands implemented. So, to keep my readme up to date, I needed to track this folder.

To setup this tracking, I will do :

```
doc-hunt add config README.md cmd
```

It will create a file .doc-hunt at the root of my repository, I have to track this file with my version control system.

Later on, after few changes on my source code, I want to check if I have to update my doc, I run :

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

I will check if I made changes in cmd/check.go that could be relevant to update my readme.

When everything is fine, I update doc-hunt doing :

```
doc-hunt update
```

And I commit all changes occuring on my file .doc-hunt.

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
