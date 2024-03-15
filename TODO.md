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
