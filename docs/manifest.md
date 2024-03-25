# Manifest



## `skip` Property

A list of files and directories (globbing is supported) to completely skip,
will not be processed or copied to the output directory

If you need to have a file be part of the template, but not renderd in the
output directory, then list it as part

## `copyAsIs` Property

Any type of file can be placed in the template, however you may not want to
avoid sending binary files through Go's template engine. `copyAsIs` allows
you to have them copied to the output without alteration.

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
