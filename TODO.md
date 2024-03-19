1. Rename the application to TmplPress
    1. rename to template.schema.json to tmplpress.schema.json
2. Feature: Validate template.json
   1. Required fields
   2. Check validation
   3. Fields that are not part of the schema.
3. Verify that copying empty directories from the substitute directory does not leave the empty file.
4. Replace cli.StringMap or stdlib.StringMap with local strMap.
5. Search to remove any reference to "zip" or "archive."
6. Remove the .git directory after cloning, then remove logic looking for .git to skip.
7. Add "validate" to the subcommand manifest. To validate a template manifest.
