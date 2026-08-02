package cmd

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/form/formdefs"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/packetver"
	"github.com/spf13/pflag"
)

var safeDir string

func Run(args []string) (err error) {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("PANIC", "p", p, "stack", debug.Stack())
			panic(p)
		}
	}()
	if len(args) == 0 {
		err = shell()
	} else {
		err = run(args)
	}
	if err != nil && err != ErrQuit {
		cio.Open().Error(err)
	}
	return err
}

func run(args []string) (err error) {
	switch args[0] {
	case "cd", "chdir":
		return cmdChdir(args[1:])
	case "connect", "c":
		return cmdConnect(args[1:])
	case "delete":
		return cmdDelete(args[1:])
	case "dump":
		return cmdDump(args[1:])
	case "edit", "e":
		return cmdEdit(args[1:])
	case "forms":
		return cmdForms(args[1:])
	case "gui", "web":
		return cmdGUI(args[1:])
	case "help", "h", "--help", "-?":
		return cmdHelp(args[1:])
	case "ics309", "309":
		return cmdICS309(args[1:])
	case "import":
		return cmdImport(args[1:])
	case "list", "l", "ls", "log":
		return cmdList(args[1:])
	case "manual", "man":
		return cmdManual(args[1:])
	case "mark":
		return cmdMark(args[1:])
	case "mkdir":
		return cmdMkdir(args[1:])
	case "n", "new":
		return cmdNew(args[1:])
	case "outpost":
		return cmdOutpost(args[1:])
	case "pdf":
		return cmdPDF(args[1:])
	case "print":
		return cmdPrint(args[1:])
	case "pwd":
		return cmdPwd(args[1:])
	case "q", "quit", "exit":
		return ErrQuit
	case "server":
		return cmdServer(args[1:])
	case "set":
		return cmdSet(args[1:])
	case "s", "show":
		return cmdShow(args[1:])
	case "version":
		return cmdVersion(args[1:])
	case "config", "files", "script":
		return cmdHelp(args[1:])
	default:
		return errors.NewF("No such command %q.", args[0])
	}
}

// TODO: handle the "safe directory" issues in the command line.
func shell() (err error) {
	c := cio.Open()
	if !c.InputIsTerm || !c.OutputIsTerm {
		return ErrUsage("usage: packet «command»\n       packet help\n")
	}
	c.Welcome(`SCCo Packet Messenger v%s.  Type "help" for help.`, packetver.Version)
	for {
		var (
			line    string
			args    []string
			in      *os.File
			out     *os.File
			saveIn  *os.File
			saveOut *os.File
		)
		// Read and parse the command line.
		if line, err = c.ReadCommand(); err != nil {
			if err == io.EOF {
				line = "quit"
			} else {
				return err
			}
		}
		if args, in, out, err = parseCommandLine(line); err != nil {
			c.Error(err)
			continue
		}
		if len(args) != 0 && args[0] == "packet" {
			// Just in case they type "packet foo" while already in
			// the shell.
			args = args[1:]
		}
		if len(args) == 0 {
			continue
		}
		// Save the old stdin and stdout and apply the new ones.
		saveIn, saveOut = os.Stdin, os.Stdout
		if in != nil {
			os.Stdin = in
		}
		if out != nil {
			os.Stdout = out
		}
		c.Detect()
		// Run the command.
		err = run(args)
		// Restore the old stdin and stdout.
		os.Stdin, os.Stdout = saveIn, saveOut
		if in != nil {
			in.Close()
		}
		if out != nil {
			out.Close()
		}
		c.Detect()
		// Handle the result of the command.
		if err != nil && err != ErrQuit {
			c.Error(err)
		}
		if err == ErrQuit {
			return nil
		}
	}
}

// parseCommandLine parses a command line that the user typed at our command
// line.  It interprets redirection.
func parseCommandLine(line string) (args []string, in, out *os.File, err error) {
	args = tokenizeLine(line)
	i := 0
	for i < len(args) {
		if args[i] == "<" {
			if in != nil {
				return nil, nil, nil, errors.New("There are multiple input redirects on the command line.")
			}
			if i == len(args)-1 || args[i+1] == "<" || args[i+1] == ">" || args[i+1] == ">>" {
				return nil, nil, nil, errors.New("There is no filename after '<'.")
			}
			if in, err = os.Open(args[i+1]); err != nil {
				return nil, nil, nil, err
			}
			args = append(args[:i], args[i+2:]...)
			continue
		}
		if args[i] == ">" || args[i] == ">>" {
			if out != nil {
				return nil, nil, nil, errors.New("There are multiple output redirects on the command line.")
			}
			if i == len(args)-1 || args[i+1] == "<" || args[i+1] == ">" || args[i+1] == ">>" {
				return nil, nil, nil, errors.New("There is no filename after '" + args[i] + "'.")
			}
			if args[i] == ">" {
				out, err = os.Create(args[i+1])
			} else {
				out, err = os.OpenFile(args[i+1], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			}
			if err != nil {
				return nil, nil, nil, err
			}
			args = append(args[:i], args[i+2:]...)
			continue
		}
		i++
	}
	return args, in, out, nil
}

// tokenizeLine parses a received command line into words.  It supports
// rudimentary quoting with either ' or ".  It does not support backslash
// escapes.  It treats unquoted <, >, and >> as separate words even when not
// surrounded by whitespace.
func tokenizeLine(line string) (args []string) {
	var partial string
	var quoted bool

	for line != "" {
		idx := strings.IndexAny(line, " \t\f\r\n'\"<>")
		if idx < 0 {
			partial += line
			args = append(args, partial)
			return args
		}
		if line[idx] == '\'' || line[idx] == '"' {
			partial += line[:idx]
			idx2 := strings.IndexByte(line[idx+1:], line[idx])
			if idx2 < 0 {
				partial += line[idx+1:]
				args = append(args, partial)
				return args
			}
			idx2 += idx + 1 // make it an offset into line
			partial += line[idx+1 : idx2]
			quoted = true
			line = line[idx2+1:]
			continue
		}
		if line[idx] == '>' && idx < len(line)-1 && line[idx+1] == '>' {
			if partial != "" || quoted {
				args, partial, quoted = append(args, partial), "", false
			}
			args, line = append(args, line[idx:idx+2]), line[idx+2:]
			continue
		}
		if line[idx] == '<' || line[idx] == '>' {
			if partial != "" || quoted {
				args, partial, quoted = append(args, partial), "", false
			}
			args, line = append(args, line[idx:idx+1]), line[idx+1:]
			continue
		}
		partial += line[:idx]
		if partial != "" || quoted {
			args, partial, quoted = append(args, partial), "", false
		}
		line = line[idx+1:]
	}
	if partial != "" || quoted {
		args = append(args, partial)
	}
	return args
}

// incRead performs a read operation on the incident in the current directory,
// if any.
func incRead(fn func(*incident.Incident) error) (err error) {
	var dir string

	if dir, err = os.Getwd(); err != nil {
		return errors.NewF("Unable to determine current directory: os.Getwd: %s", err)
	}
	if dir, err = filepath.EvalSymlinks(dir); err != nil {
		return errors.NewF("Unable to resolve current directory: filepath.EvalSymLinks: %s", err)
	}
	if !incident.IsIncident(dir) {
		return errors.NewF("The current directory %s is not an incident directory.", dir)
	}
	return incident.Read(dir, fn)
}

func incWrite(create bool, fn func(*incident.Incident) error) (err error) {
	var dir string

	if dir, err = os.Getwd(); err != nil {
		return errors.NewF("Unable to determine current directory: os.Getwd: %s", err)
	}
	if dir, err = filepath.EvalSymlinks(dir); err != nil {
		return errors.NewF("Unable to resolve current directory: filepath.EvalSymLinks: %s", err)
	}
	if incident.IsIncident(dir) {
		return incident.Write(dir, fn)
	}
	if !create {
		return errors.NewF("The current directory %s is not an incident directory.", dir)
	}
	if incident.IsUnsafeIncidentDir(dir) {
		return errors.NewF("The current directory %s is not a proper directory to create an incident in.  Make a subdirectory with an incident-specific name instead.", dir)
	}
	return incident.Create(dir, fn)
}

func gaveMutuallyExclusiveFlags(set *pflag.FlagSet, flags ...string) (err error) {
	var seen bool

	for _, flag := range flags {
		if f := set.Lookup(flag); f.Changed {
			if seen {
				if len(flags) == 2 {
					return errors.NewF("The --%s and --%s flags are incompatible.", flags[0], flags[1])
				} else {
					return errors.NewF("The --%s and --%s flags are incompatible.", strings.Join(flags[:len(flags)-1], ", --"), flags[len(flags)-1])
				}
			}
			seen = true
		}
	}
	return nil
}

var registerFormsOnce sync.Once

func registerForms() {
	registerFormsOnce.Do(func() {
		var c = cio.Open()

		if err := formdefs.CheckForUpdates(false, true); err != nil {
			c.Warn("%s", err)
		} else if c.OutputIsTerm {
			if readme := formdefs.GetFlushREADME(); readme != "" {
				c.Confirm("***** NEW FORMS INSTALLED *****\n%s\n*******************************\n", readme)
			}
		}
		if err := formdefs.RegisterForms(); err != nil {
			c.Warn("%s", err)
		}
	})
}
