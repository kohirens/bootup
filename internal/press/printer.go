package press

import (
	"github.com/kohirens/stdlib/cli"
	"os"
)

const (
	dirMode      = 0744
	emptyFile    = ".empty"
	gitConfigDir = ".git"
	maxTmplSize  = 1e+7
	PS           = string(os.PathSeparator)
)

type AnswersJson struct {
	Placeholders cli.StringMap `json:"placeholders"`
}
