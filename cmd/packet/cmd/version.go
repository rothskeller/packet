package cmd

const (
	versionSlug = `Prints the current software version`
	versionHelp = `
usage: packet version

The "packet version" command displays the current software version number.  (For the versions of installed forms, use the "packet forms list" command.)
`
)

func cmdVersion(args []string) (err error) {
	panic("not implemented")
}
