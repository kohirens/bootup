1. WIP: Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
2. Refactor FEC to the template.json schema as ignoreExtension list.
3. Rename "excludes" in the template.json as "copy" as is. Files in this list
   will NOT be processed through the template engine and copied-as-is.
4. When generating a template.json manifest, if there is already a template.json, load it, then merge it with new data
5. Rename the application to TmplPress
    1. rename to template.schema.json to tmplpress.schema.json
