1. Added "substitute" and "substituteDir" to replace "replace" feature.
   1. Add all files to the skip array.
   2. Avoid adding the substitute directory to the skip array and allow the
      files in the substitute directory to be processed, this is by default.
   3. Do not replace files until the template processing completes.
   4. Then move all files out of the substituteDir into the root. The files
      should not exist so there should be no need to force move.
   5. Change "substitute" to a string instead of an object. We do not need the
      map to get a list of the files in the directory. And we will not allow
      renaming at this time.
2. Add basic (bool, int, unsigned) validation for placeholder.
3. Remove any zip or 7zip extract support (git support only)
4. WIP: Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
5. Rename the application to TmplPress
   1. rename to template.schema.json to tmplpress.schema.json
6. Refactor FEC to the template.json schema as ignoreExtension list.
7. Rename "excludes" in the template.json as "copy" as is. Files in this list
   will NOT be processed through the template engine and copied-as-is.
