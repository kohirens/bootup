# TmplPress

A Go Template Rendering Tool. Start a project from a template rather than from scratch.

## Table of Contents

* [TmplPress](#tmplpress)
    * [Info](#info)
    * [Description](#description)
    * [What Is A Template](#what-is-a-template)
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

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/tmplpress/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/tmplpress/tree/main)

## Description

Designed for initializing applications; it can be applied generally. For
example, configuration of a dynamic CI/CD pipeline (very useful). Don't limit
your imagination and, use it for any type of project where you need to fill-in
values in multiple files within a single directory. Initialize those files
with specific values, especially in an automated repeatable way.
Where you can benefit from saving time and reducing errors.

## What Is A Template

Template is the term used for the concept and actual __template__ file. Any
folder containing a **manifest*** and one or more files that contain
[Go template Actions] markup count as a template. Typically, you can point to
a Git repository or a local folder will suffice.

* The manifest is a configuration file, in JSON format, who's properties provide
details to help __**"press"**__ (or render) the template. The term press as in
newspaper press.

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
go install github.com/kohirens/tmplpress
```

### Using Docker

```
docker pull kohirens/tmplpress:x.x.x
```

### Using Pre-built Binary

```
mkdir -p "${HOME}/bin"
curl -L -o tmplpress.tar.gz https://github.com/kohirens/tmplpress/releases/download/x.x.x/tmplpress-linux-amd64.tar.gz
tar -xzvf tmplpress.tar.gz  ${HOME}/bin
export PATH="${HOME}/bin:${PATH}"
```

## Using a Template

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
[Go template Actions]: https://pkg.go.dev/text/template#hdr-Actions
