# How To Build A Template JSON Manifest

Your template will need a manifest `template.json` file. It is safe to build
one manually by using the following format:

```JSON
{
    "version": "2.2.0",
    "emptyDirFile": ".empty",
    "placeholders": {
        "appName": "a name for the application",
        "repoName": "a repository name for the application",
        "repoOrg": "your GitHub organization name"
    }
}
```

Or generate one using the `tmplpress manifest generate` command,
supplying the path to the template as the argument. The generated file will be
placed in that path supplied. NOTE: If a manifest alread exist, it will
be updated to:
1. A new format based on the version of schema that `tmplpress` supports.
2. Updated placeholders to reflect any added/removed.

At minimum the `template.json` needs to contain

1. A `version` property with the desired template.json schema version.
2. A `emptyDirFile` property with the name of file that represents a directory
   as empty.

## Placeholders

Placeholders are template actions that take a value to replace the variable.
Placeholders such at `{{ .AppName }}` will correlate to placeholders in the
`template.json` file:

```json
{
    "version": "2.2.0",
    "emptyDirFile": ".empty",
    "placeholders": {
        "AppName": "a name for the application",
    }
}
```

All variables must be provided a value, even if that is the empty string.
Notice the string values for each key in the `placeholders` property equates
to a question. This is because they can be used as prompts to
ask for the value when filing out the template from the CLI.

## References

* [JSON Schema](https://json-schema.org/learn/getting-started-step-by-step#intro)
