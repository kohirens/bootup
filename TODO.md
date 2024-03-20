1. Feature: Validate template.json
    1. Required fields
    2. Check validation
    3. Fields that are not part of the schema.
    4. Add "validate" to the subcommand manifest. To validate a template manifest.
2. Verify that copying empty directories from the substitute directory does not leave the empty file.
3. Search to remove any reference to "zip" or "archive."
4. Remove the .git directory after cloning, then remove logic looking for .git to skip.
