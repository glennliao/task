package tasker

import (
	"github.com/glennliao/task/tasker/util"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

var logger = util.Logger{}

func (t *Tasker) Init(taskerJs string) {
	dir, _ := homedir.Dir()

	var ConfigRoot = filepath.Join(dir, ConfigParentDir)
	t.configRoot = filepath.Join(ConfigRoot, ConfigDir)
	t.configUseRoot = filepath.Join(t.configRoot, ConfigUseDir)

	os.Mkdir(ConfigRoot, os.ModePerm)
	os.Mkdir(t.configRoot, os.ModePerm)

	t.taskerJs = taskerJs

	f, err := os.Create(filepath.Join(t.configRoot, ConfigTaskerJsFilename))
	if err != nil {
		panic(err)
	}
	f.WriteString(taskerJs)
	f.Close()
}
