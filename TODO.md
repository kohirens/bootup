1. When generating a template.json manifest, if there is already a template.json, load it, then merge it with new data
2. Rename the application to TmplPress
    1. rename to template.schema.json to tmplpress.schema.json
3. Feature: Validate template.json
   1. Required fields
   2. Check validation
   3. Fields that are not part of the schema.
4. Verify that copying empty directories from the substitute directory does not leave the empty file.
5. Replace cli.StringMap or stdlib.StringMap with local strMap.
6. Search to remove any reference to "zip" or "archive."
7. Remove the .git directory after cloning, then remove logic looking for .git to skip.
8. Add "validate" to the subcommand manifest. To validate a template manifest.
