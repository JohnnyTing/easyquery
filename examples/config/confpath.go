package config

import (
	"path/filepath"

	"github.com/ztrue/tracerr"

	"github.com/integralist/go-findroot/find"
)

var (
	ConfPath string
)

func init() {
	root, err := find.Repo()
	if err != nil {
		err = tracerr.Errorf("find repo Error: %s", err.Error())
		panic(err)
	}
	ConfPath = root.Path + string(filepath.Separator) + "examples" + string(filepath.Separator) + "config" + string(filepath.Separator)
}
