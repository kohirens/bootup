1. Rename "excludes" in the template.json as "copyAsIs". Files in this list
   will NOT be processed through the template engine and copied-as-is.
2. When generating a template.json manifest, if there is already a template.json, load it, then merge it with new data
3. Rename the application to TmplPress
    1. rename to template.schema.json to tmplpress.schema.json
4. Feature: Validate template.json
   1. Required fields
   2. Check validation
   3. Fields that are not part of the schema.
5. Verify that copying empty directories from the substitute directory does not leave the empty file.
6. Replace cli.StringMap or stdlib.StringMap with local strMap.
7. Search to remove any reference to "zip" or "archive."
8. Remove the .git directory after cloning, then remove logic looking for .git to skip.
9. Change the manifest command to require the word "generate" in order to
   generate a new template.json.
10. Add "validate" to the subcommand manifest. To validate a template manifest.
