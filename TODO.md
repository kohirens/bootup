1. Refactor FEC to the template.json schema as ignoreExtension list.
2. Rename "excludes" in the template.json as "copy" as is. Files in this list
   will NOT be processed through the template engine and copied-as-is.
3. When generating a template.json manifest, if there is already a template.json, load it, then merge it with new data
4. Rename the application to TmplPress
    1. rename to template.schema.json to tmplpress.schema.json
5. Feature: Validate template.json
   1. Required fields
   2. Check validation
   3. Fields that are not part of the schema.
