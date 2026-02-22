package cio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

// EditResult is the result of an EditField operation.
type EditResult byte

// Values for EditResult:
const (
	// ResultEOF says that the edit was terminated by reaching end of file
	// on standard input.  (This only happens in batch mode.)
	ResultEOF EditResult = '\004'
	// ResultNext says that the editor should move on to the next field.
	ResultNext EditResult = 'v'
	// ResultPrevious says that the edit should move back to the previous
	// field.
	ResultPrevious EditResult = '^'
	// ResultDone says that the edit should stop.
	ResultDone EditResult = '\033'
)

const editorHelp = `Editor: [F1]=Help [Tab]=Next [Shift-Tab]=Prev [ESC]=Exit [Ctrl-C]=Abort`

// StartEdit prints the editor help text that appears at the beginning of an
// editor session.
func (cio *CIO) StartEdit() {
	if !cio.InputIsTerm || !cio.OutputIsTerm {
		return
	}
	cio.print(colorHelp, setLength(editorHelp, cio.Width-1))
	cio.print(0, "\n")
}

type modefunc func() (modefunc, EditResult, error)

type editor struct {
	cio        *CIO
	label      string
	labelWidth int
	value      string
	valueWidth int
	choices    []string
	help       string
	hint       string
	multiline  bool
	obscured   bool
	cursor     int
	sels, sele int
	changed    bool
}

// EditField gets a new value for a field.  label is the field label.
// labelWidth is the width of the label column (>=len(label)).  value is the
// current value of the field.  valueWidth is the width of the value entry
// field, which should map to the amount of space available to render the field
// in its PDF layout.  choices is an optional list of suggested values for the
// field.  help is the help text for the field.  hint is the text that appears
// in the entry area when the field has no value.  multiline is a flag saying
// that a multi-line value for the field is expected.  obscured is a flag saying
// that the field contents should not be shown (e.g., a password field).
// cleanup is an optional function that takes the entered value and applies
// autocorrections to it (e.g., capitalizing an all-caps field, correcting date
// punctuation, etc.).  The function returns the new value of the field and an
// indication of what caused the edit to finish.  It returns an non-nil error if
// the edit was aborted and the new value should not be used.
func (cio *CIO) EditField(
	label string, labelWidth int, value string, valueWidth int,
	choices []string, help, hint string, multiline, obscured bool,
	cleanup func(string) string,
) (result EditResult, newvalue string, err error) {
	if !cio.InputIsTerm || !cio.OutputIsTerm {
		// Not interactive, so just read the new value from stdin.
		return readFieldStdin()
	}
	cio.rawMode()
	defer cio.restoreTerminal()
	e := editor{
		cio:        cio,
		label:      label,
		labelWidth: max(labelWidth, len(label)),
		value:      value,
		valueWidth: valueWidth, // modified below
		choices:    choices,
		help:       help,
		hint:       hint,
		multiline:  multiline,
		obscured:   obscured,
		cursor:     len(value),
		sele:       len(value),
	}
	for _, c := range e.choices {
		e.valueWidth = max(e.valueWidth, len(c))
	}
	var mode modefunc
	if len(e.choices) != 0 && (e.value == "" || slices.Contains(e.choices, e.value)) {
		mode = e.choicesMode
	} else if labelWidth+2+len(e.value) >= cio.Width || strings.Contains(e.value, "\n") {
		mode = e.multilineMode
	} else {
		mode = e.onelineMode
	}
	for mode != nil && err == nil {
		mode, result, err = mode()
	}
	if err != nil {
		return 0, "", err
	}
	if cleanup != nil {
		e.value = cleanup(e.value)
	}
	e.displayResult()
	return result, e.value, nil
}

// displayResult removes the editor from the screen and replaces it with the
// result of the editing.
func (e *editor) displayResult() {
	var (
		value   string
		linelen int
		lines   []string
		indent  string
	)
	e.cio.print(colorLabel, e.label)
	value = strings.TrimRight(e.value, "\n")
	if e.obscured {
		value = obscureValue(value)
	}
	// Find the length of the longest line in the value.
	lines = strings.Split(value, "\n")
	for _, line := range lines {
		linelen = max(linelen, len(line))
	}
	if linelen <= e.cio.Width-e.labelWidth-3 {
		e.cio.print(0, spaces[:e.labelWidth+2-len(e.label)])
		indent = spaces[:e.labelWidth+2]
	} else {
		lines, _ = wrap(value, e.cio.Width-5)
		e.cio.print(0, "\n    ")
		indent = spaces[:4]
	}
	for i, line := range lines {
		if i != 0 {
			e.cio.print(0, indent)
		}
		e.cio.print(0, line)
		e.cio.print(0, "\n")
	}
}

// showHelp displays the help text for the field.
func (e *editor) showHelp() {
	e.cio.move(0, 0)
	e.cio.cleanTerminal()
	lines, _ := wrap(e.help, e.cio.Width-1)
	for _, line := range lines {
		e.cio.print(colorHelp, setLength(line, e.cio.Width-1))
		e.cio.print(0, "\n")
	}
	e.cio.print(colorHelp, setLength(editorHelp, e.cio.Width-1))
	e.cio.print(0, "\n")
	e.cio.cleanTerminal()
}

// readFieldStdin reads the value of a field from standard input (i.e.,
// EditField when running non-interactively).
func readFieldStdin() (EditResult, string, error) {
	contents, err := io.ReadAll(os.Stdin)
	if err != nil {
		return 0, "", fmt.Errorf("reading stdin: %s", err)
	}
	contents = bytes.TrimRight(contents, "\r\n")
	contents = bytes.ReplaceAll(contents, []byte{'\r', '\n'}, []byte{'\n'})
	contents = bytes.ReplaceAll(contents, []byte{'\r'}, []byte{'\n'})
	return ResultEOF, string(contents), nil
}

// obscureValue returns a string of the same length as the input, consisting
// entirely of stars.  It is used for password fields.
func obscureValue(s string) string {
	return strings.Repeat("*", len(s))
}
