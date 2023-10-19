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
func Substitute(dir, tp string) error {
	if !path.Exist(dir) {
		return nil
	}

	if !path.Exist(tp) {
		return fmt.Errorf(msg.Stderr.PathNotExist, tp)
	}

	// This has to overwrite files, delete the direcory would be bad as
	// there may be files that the template designer place there that was meant
	// to keep.
	if e := path.CopyDirToDir(dir, tp, PS, dirMode); e != nil {
		return fmt.Errorf(msg.Stderr.CannotCopyDirToDir, dir, tp, e.Error())
	}

	return nil
}
