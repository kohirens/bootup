package press

import (
	"fmt"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/msg"
)

// Substitute Looks for a special directory, can be specified in the
// template.json, then copies all files from that directory to the root of
// the TmplPath, overwriting any existing files.
// This is done before any templates are processed, to ensure they are run
// also through the template engine.
func Substitute(source, dest string) error {
	// This has to overwrite files, delete the directory would be bad as
	// there may be files that the template designer place there that was meant
	// to keep.
	if e := path.CopyDirToDir(source, dest, PS, dirMode); e != nil {
		return fmt.Errorf(msg.Stderr.CannotCopyDirToDir, source, dest, e.Error())
	}

	return nil
}
