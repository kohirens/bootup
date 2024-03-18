# TmplToApp

Start an app (or something) from a template.

## Table of Contents

* [TmplToApp](#tmpltoapp)
    * [Info](#info)
    * [Description](#description)
    * [Installation](#installation)
        * [Requirements](#requirements)
        * [Using Go](#using-go)
        * [Using Docker](#using-docker)
        * [Using cURL & tar](#using-curl--tar)
    * [Using a Template](#using-a-template)
    * [Making a Template](/docs/template-designing.md#making-a-template)
    * [FYI](#fyi)
    * [Teminology](/docs/terminology.md)

## Info

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/tmpltoapp/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/tmpltoapp/tree/main)

## Description

A template is a collection of files organized in a folder hierarchy. Any
extension can be used as long as it is text (only tested with UTF-8) containing
Go template syntax. This application takes such a folder and processes each
file in the folder structure to an output folder of your choosing.

Data for the template is supplied with questions answer from the CLI or a
JSON file as input. This is extremely powerful; you are only limited to your
knowledge of Go templates.

You can make whole project templates or smaller pieces your more likely to use
on a regular basis. For examole, making a Docker file template or a CI/CD
configuration you use for many projects. Making a tempalte out of them to fill
in application details for example.

The idea is to quickly setup things you need on a regular basis.

**Hint:** Templates are invaluable for quickly setting up apps/projects layouts
(even a small parts) that you commonly use. This is especially true when using
the `answer.json` file with automation.

## Installation

### Requirements

Git must be installed on your system in order to use this tool. Git is used
to perform actions such as cloning and checking out branches or tags, and is
necessary for this application to perform its functions.

### Using Go

```
go install github.com/kohirens/tmpltoapp
```

### Using Docker

```
docker pull kohirens/tmpltoapp:x.x.x
```

### Using Pre-built Binary

```
mkdir -p "${HOME}/bin"
curl -L -o tmpltoapp.tar.gz https://github.com/kohirens/tmpltoapp/releases/download/x.x.x/tmpltoapp-linux-amd64.tar.gz
tar -xzvf tmpltoapp.tar.gz  ${HOME}/bin
export PATH="${HOME}/bin:${PATH}"
```

## Using a Template

You'll need to download this tool in order to use a template. See [Installation]
if you have not done so.

NOTE: There are command line flags should you need to place the arguments
out of order. Run the program with `-h` or `--help` for options.

Run this application with 3 parameters:
1. a path to a template, a URL or local folder.
2. a path to where you want to place the project (it should not exist).

### Notes About Template Processing

* All variables are treated as strings.
* If any variables are in the `template.json` that are supplied by an answer JSON, then processing will halt and ask for them. 
* Empty directories will be placed without the ".empty" file.
* Files listed in the `excludes` list are output to the final app directory without template processing.
* Template are processed with the Go lib [Golang text/template].

## FYI

1. Why is it called "TmplPress"?
    * The name is a play on newspress. Old machines used to print newspapers.
      Like a newspaper the TmplPress (Template Press) produces copies from
      templates.
2. What is up with the name "printer.go" in the press package"
    * Going along with the them of newspress, the machine as a whole acts as a
      printer. Originally the name was lever, for what a person would pull to
      print 1 side of a newpaper, but since the press package is meant to
      contain all the parts build to produce the paper, it made more since to
      call the file that contains the main function to produce a tempalte be
      name printer.

---

[Golang text/template]: https://golang.org/pkg/text/template/
