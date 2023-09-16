# TmplToApp

Start an app (or something) from a template.

## Table of Contents

* [TmplToApp](#tmpltoapp)
    * [Description](#description)
    * [Installation](#installation)
        * [Using Go](#using-go)
        * [Using Docker](#using-docker)
    * [Get Started](#get-started)
        * [Making a Template](#making-a-template)
        * [Using a Template](#using-a-template)


## Info

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/tmpltoapp/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/tmpltoapp/tree/main)

## Description

A template is a collection of files organized in a folder hierarchy. Any
extension can be used, but the template must be written in text (only tested
with UTF-8) containing Go template syntax. This application takes such a folder
and processes each file in the folder structure to an output folder of your
choosing.

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

## Terminology

* Placeholders - refer to the variables that need to be filled in when a file
  is process as a Go template.
* Template - refer to top/root folder/directory as a whole, which contains
  text files (of any extension) containing Go template syntax.
* Empty directory - a directory with a single file named `.empty`, contents
  are ignored.
* Templates source - the Git repository for the __Template__ repo.

## Installation

### Using Go

```
go install github.com/kohirens/tmpltoapp
```

### Using Docker

```
docker pull kohirens/tmpltoapp:x.x.x
```

### Using cURL on Linux

```
cd /tmp
mkdir -p "${HOME}/bin"
curl -L -o tmpltoapp.tar.gz https://github.com/kohirens/tmpltoapp/releases/download/x.x.x/tmpltoapp-linux-amd64.tar.gz
tar -xzvf tmpltoapp.tar.gz  ${HOME}/bin
export PATH="${HOME}/bin:${PATH}"
```

## Get Started

### Making a Template

You can quickly make a Template by following these steps.

1. Make a new directory with a name of your choosing.
2. Make a README.md and add `# {{.AppName}}` as the content.
3. Open a command line to this folder and run `tmpltoapp manifest ./`. This
   will generate the manifest `template.json` file containing some details about
   your template. Mainly the placeholder.
4. You can edit this file by giving the AppName key a value like:
   "application name". This acts as a label or question when someone uses your
   template. More on that later.
5. Run `git inti` and then `git add .`, then commit the changes.

That is the start of your template. But you can add more folders,
and if a directory should be empty, then place a file named ".empty" in it.
The .empty file does need any text in it.

Add more files as needed, the extension does not matter as it will be
processed as a Go template, unless it is excluded in the template.json manifest.
See [How To Build A Template JSON Manifest] for other details that can be added.

### Using a Template

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

---

[Golang text/template]: https://golang.org/pkg/text/template/
[How To Build A Template JSON Manifest]: /docs/building-a-template-json.md
