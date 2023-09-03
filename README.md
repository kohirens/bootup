* [TmplToApp](#tmpltoapp)
    * [Description](#description)
    * [Installation](#installation)
        * [Using Go](#using-go)
        * [Using Docker](#using-docker)
    * [Get Started](#get-started)
        * [Making a Template](#making-a-template)
        * [Using a Template](#using-a-template)

# TmplToApp

Start an app (or something) from a template.

## Info

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/tmpltoapp/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/tmpltoapp/tree/main)

## Description

Uses the Go text/template engine to produce output such as an app, project, etc.
A template is a folder with files (of any extension). Each file processed as a
Go template. Data for the template is supplied with JSON files. They can be
simple key/value pairs. However, you are only limited to your knowledge of Go
templates. Meaning it can be more than just simple string replacement of
placeholders. Think loops, conditions, and function calls (built into go).

This can be used to process anything you can start from a template. For example
image you want to send out an email to many clients, then you can process
a file, in a loop, and supply the data that changes to produce different output
with each pass of the loop.

## Terminology

* Placeholders - refer to the variables that need to be filled in when a file
  is process as a Go template.
* Template - is a folder/directory with files (of any extension).
* Empty directory - a directory with a single file named `.empty`, contents
  are ignored.
* Templates source - can be a local folder, URL to a zip file, or Git repo.

## Installation

### Using Go

```
go get github.com/kohirens/tmpltoapp
```

### Using Docker

```
docker pull kohirens/tmpltoapp:latest
```

### Using cURL

```
cd /tmp
mkdir -p "${HOME}/bin"
curl -L -o tmpltoapp.tar.gz https://github.com/kohirens/tmpltoapp/releases/download/x.x.x/tmpltoapp-linux-amd64.tar.gz
tar -xzvf tmpltoapp.tar.gz  ${HOME}/bin
export PATH="${HOME}/bin:${PATH}"
```

## Get Started

### Making a Template

Templates are just directories containing files, which can contain Go template
syntax that will be processed to fill in placeholders. The output will serve as
a new project.

1. Make a new directory.
2. Add folders, and if a directory should be empty, then place a file named
   ".empty" in it. It does need any text in it.
3. Add your files, the extension does not matter as it will be processed as a Go template. For example imagine "README.md" with the content.
   1. Files can contain a GoLang template placeholder, so `README.md` can contain:
      ```gotemplate
      # {{ .appName }}

      ## Summary
      ```
   Note: the placeholder `{{ .appName }}` will be replaced with the apps name at runtime.
4. Add a `template.json` file that serves as a manifest of all variables in the template, see [How To Build A Template JSON Manifest](/docs/building-a-template-json.md)
5. Commit the changes and push up to your repo.

### Using a Template

Run this application with 3 parameters:
1. a path to a template, a URL to a zip or local folder.
2. a path to where you want to place the project (it should not exist).
3. a path to an answer (JSON) file containing key/value pairs that will
   serve as variables. The name is not imp For example, an `{{ .author }}`
   placeholder would take a file that has:
   ```
   {
      "author": "your name here",
      "appName": "awesomeAppName"
   }
   ```
NOTICE: Name change pending... "start-from-tmpl" to make it obvious what this
tool does.

4. make a `template.json` file which acts as a manifest, it needs the following:
   ```
   {
      "version": "1.0.0",
      "placeholders": {
        "appName": "name of the app",
        "repoName": "name of the VCS repository"
      },
      "excludes": [
         ".gitignore",
         ".gif",
         ".jpg",
         ".png",
         ".mov",
      ]
   }
   ```
NOTE: There are command line flags should you need to place the arguments
      out of order. Run the program with `-h` or `--help` for options.

### Notes About Template Processing

* All variables are treated as strings.
* If any variables are in the `template.json` that are supplied by an answer JSON, then processing will halt and ask for them. 
* Empty directories will be placed without the ".empty" file.
* Files listed in the `excludes` list are output to the final app directory without template processing.
* Template are processed with the Go lib [Golang text/template].

---

[Golang text/template]: https://golang.org/pkg/text/template/
