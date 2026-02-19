package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
)

const helpSlug = `Print help for packet commands or topics`
const topHelp = `
The "packet" command provides multiple commands for handling packet radio messages.  When invoked with a command on the command line, it runs that command.  When invoked without any arguments, it starts a shell that allows running multiple commands without the "packet" prefix on each.

Available commands include:
  cd       ⇥` + chdirSlug + `
  connect  ⇥` + connectSlug + `
  delete   ⇥` + deleteSlug + `
  dump     ⇥` + dumpSlug + `
  edit     ⇥` + editSlug + `
  forms    ⇥` + formsSlug + `
  gui      ⇥` + guiSlug + `
  help     ⇥` + helpSlug + `
  ics309   ⇥` + ics309Slug + `
  list     ⇥` + listSlug + `
  manual   ⇥` + manualSlug + `
  mark     ⇥` + markSlug + `
  mkdir    ⇥` + mkdirSlug + `
  new      ⇥` + newSlug + `
  outpost  ⇥` + outpostSlug + `
  pdf      ⇥` + pdfSlug + `
  print    ⇥` + printSlug + `
  pwd      ⇥` + pwdSlug + `
  quit     ⇥` + quitSlug + `
  server   ⇥` + serverSlug + `
  set      ⇥` + setSlug + `
  show     ⇥` + showSlug + `
  version  ⇥` + versionSlug + `
For help on a command, run "packet help «command»".

Additional help is available on the following topics:
  config   ⇥` + configSlug + `
  files    ⇥` + filesSlug + `
  params   ⇥` + paramsSlug + `
  script   ⇥` + scriptSlug + `
For these topics, run "packet help «topic»".

This SCCo Packet software was written by Steve Roth KC6RSC.  Source code, licensing details, and bug tracker are at github.com/rothskeller/packet.
`

func cmdHelp(args []string) (err error) {
	var (
		helpText string
		c        = cio.Open()
	)
	if len(args) != 0 {
		switch args[0] {
		case "cd", "chdir":
			helpText = chdirHelp
		case "config":
			helpText = configHelp
		case "connect", "c":
			helpText = connectHelp
		case "delete":
			helpText = deleteHelp
		case "dump":
			helpText = dumpHelp
		case "edit", "e":
			helpText = editHelp
		case "files":
			helpText = filesHelp
		case "forms":
			if len(args) > 1 {
				return cmdFormsHelp(args[1:])
			}
			helpText = formsHelp
		case "gui", "web":
			helpText = guiHelp
		case "ics309", "309":
			helpText = ics309Help
		case "list", "l", "ls", "log":
			helpText = listHelp
		case "manual", "man":
			if len(args) > 1 {
				return cmdManualHelp(args[1:])
			}
			helpText = manualHelp
		case "mark":
			helpText = markHelp
		case "mkdir":
			helpText = mkdirHelp
		case "n", "new":
			helpText = newHelp
		case "outpost":
			if len(args) > 1 {
				return cmdOutpostHelp(args[1:])
			}
			helpText = outpostHelp
		case "params":
			helpText = paramsHelp
		case "pdf":
			helpText = pdfHelp
		case "print":
			helpText = printHelp
		case "pwd":
			helpText = pwdHelp
		case "q", "quit", "exit":
			helpText = quitHelp
		case "script":
			helpText = scriptHelp
		case "server":
			if len(args) > 1 {
				return cmdServerHelp(args[1:])
			}
			helpText = serverHelp
		case "set":
			helpText = setHelp
		case "s", "show":
			helpText = showHelp
		case "version":
			helpText = versionHelp
		default:
			c.Error("There is no command or help topic %q.", args[0])
		}
	}
	if helpText == "" {
		helpText = topHelp
	}
	helpText = strings.TrimLeft(helpText, "\n") // Allows newline after `
	io.WriteString(os.Stdout, c.WrapText(helpText))
	return nil
}

func usage(help string) error {
	help = strings.TrimLeft(help, "\n")
	if idx := strings.Index(help, "\n\n"); idx > 0 {
		help = help[:idx]
	}
	return ErrUsage(help)
}

type ErrUsage string

func (e ErrUsage) Error() string { return string(e) }
