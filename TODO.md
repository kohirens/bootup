1. Added "substitute" and "substituteDir" to replace "replace" feature.
   1. Add all files to the skip array.
   2. Avoid adding the substitute directory to the skip array and allow the
      files in the substitute directory to be processed, this is by default.
   3. Do not replace files until the template processing completes end, then 
      just moved all files out of the substituteDir into the root. The files
      should not exist so there should be no need to force move.
2. Add basic (bool, int, unsigned) validation for placeholder.
3. WIP: Add 7zip extract support.
    * Make use of 7zip extractor in main.
4. WIP: Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
5. Rename to tmplPress
