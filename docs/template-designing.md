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

## FYI

So there can be some confusing concepts with making templates due how the
author of the program envisions how things should work and what a template
author expects.

I the author of the program will try to explain what I was thinking
in a pragmatic way (instead of adventure/story mode style that I prefer).

### Why Start Another Template Processor

I often need to start new projects and what better way to do that than start
it from a template, like some .Net or Java projects do. So since I was
learning Go I started looking for a template processing engine that would suit
my needs to spin up Go projects. To spare you the long details after looking
for a short time and experimenting with a few I decided to write my own.

I foolishly thought something like this was simple and would not take a lot of
time. Years later here we are.

### Possibly Confusing Concepts

The template.json Schema has 3 list you can add that my see way too similar. I'm
talking about `copy`, `skip`, and `substitute`. Though they've had other names
or different implementations at times. Let me try to explain the difference.

`copy` - Copies a file as it is in the template. It is copied directly to the
out directory with no modification. This is meant for binary files like images
or some text files, like LICENSE.md.

`skip` - Meant to ignore files and skip them entirely. Why would you want to do
that you ask. Simple. It is meant for files that are NOT to be added to the out
directory or processed at all. Use it when you have files that need to exist in
the template but are not part of the template. Though there is rare it is
necessary. Such as in the case that you want to add a logo for your template, or
other marketing resources. Or even files that are meant to serve for automation.
This is a CLI tool after all. You can use the templates to automate deploying
projects.

`substitute` - Use this to substitute directories/files in place of others in
the out directory. This is a string and NOT a list like `copy` and `skip`. Its
value is the name of a directory in your template (other that "substitute") that
will server as the substitution directory. Files passed through the template
processor and copied to the root of the out directory, mirroring their relative
paths.
Its almost like combining `skip` and `copy` and then allows copied files to
also go through the template processor.

Let walk through an example use case. Say you want to test your template
through an automated build system like with CircleCI. So you add a
".circleci/config.yml" configuration to do that.
However, you do NOT want that directory to show up in the out directory.
So you add it to the `skip` array of files to ignore. But what if you did want
to add a different ".circleci/config.yml" to the out directory.
Then you can add that to the copy list, however that will not work unless you
add it to your template with a different name, but then it would have the wrong
name in the final output. Plus it would not get parsed by the template
processor.
That is when the `substitute` directory comes in handy. You copy the templated
".circleci" directory it into the "substitute" directory, organizing it as you
will want it to be in the out directory. Giving you a CircleCi config that test
your template with automation, and another config that is for the user
of your template.

### Missing Features

These are features that were thought of but had no reason to implemented because
they weren't used during development and conceptualizing.

* You cannot rename files to the out directory, filenames mirror
  their relative template path.
* There is no globing. You can oly use relative directory and file names only in
  the template.json manifest. Copy allows things "*.jpg", but its really only
  the extension its looking for.
---

[How To Build A Template JSON Manifest]: /docs/building-a-template-json.md
