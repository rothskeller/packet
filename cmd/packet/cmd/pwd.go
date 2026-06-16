package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/spf13/pflag"
)

const (
	pwdSlug = `Print current directory`
	pwdHelp = `
usage: pwd

The "pwd" command displays the current directory path.
`
)

func cmdPwd(args []string) (err error) {
	flags := pflag.NewFlagSet("pwd", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"pwd"})
	} else if err != nil {
		cio.Open().Error(err)
		return usage(pwdHelp)
	}
	if len(args) != 0 {
		return usage(pwdHelp)
	}
	var dir string
	if dir, err = os.Getwd(); err != nil {
		return err
	}
	if nd, err := filepath.EvalSymlinks(dir); err == nil {
		dir = nd
	}
	fmt.Println(dir)
	return nil
}
