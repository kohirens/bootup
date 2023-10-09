## Template Designing

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

---

[How To Build A Template JSON Manifest]: /docs/building-a-template-json.md
