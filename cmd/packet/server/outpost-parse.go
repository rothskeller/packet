package server

import (
	"regexp"
	"strings"
)

var eol = "\r\n"

type ParsedEmail struct {
	Headers      map[string]string
	Fields       map[string]string
	Urgent       bool
	AddonName    string
	FormType     string
	AddonVersion string
}

// ParseEmail an email message.  It doesn't have to contain a form.
func ParseEmail(inputMessage string) (result *ParsedEmail, repaired string) {
	// Here we're using PackItForms' algorithm for parsing the message,
	// because I don't want to risk incompatibilities with parsing details
	// right now.
	result = &ParsedEmail{Fields: make(map[string]string)}
	if len(inputMessage) == 0 {
		return result, inputMessage
	}
	// Convert all line endings to CRLF.
	var email = toEOL(inputMessage)
	// Split the message into a headers string and a body string.
	var headers = ""
	var body = email
	var found = strings.Index(body, eol+eol)
	if found >= 0 {
		// This empty line separates the headers from the body.
		headers = body[:found]
		body = body[found+2*len(eol):]
	}
	if match := re(`(^|\r\n)[!#]`).FindStringSubmatchIndex(headers); match != nil {
		// A line starting with ! or # can't be a header.  Maybe there
		// are no headers.
		headers = email[:match[0]]
		body = email[match[0]+match[3]-match[2]:]
	}
	// Parse the headers into result.headers.
	var fieldName, fieldValue string
	for _, line := range strings.Split(headers, eol) {
		if fieldName != "" && re(`^\s`).MatchString(line) {
			// A continuation of the previous line.
			fieldValue += line
		} else {
			if match := re(`^(\S+)\s*:\s*`).FindStringSubmatch(line); match != nil {
				if fieldName != "" {
					result.Headers[fieldName] = fieldValue
				}
				fieldName = strings.ToLower(match[1])
				fieldValue = line[len(match[0]):]
			}
		}
	}
	if fieldName != "" {
		result.Headers[fieldName] = fieldValue
	}
	// Assume the body contains a form.  Repair any damage JNOS did to it.
	message := repairMessage(body)
	// Parse it into result.addonName, addonVersion, formType, and fields.
	fieldName, fieldValue = "", ""
	for _, line := range strings.Split(message, eol) {
		if fieldName == "" {
			if strings.HasPrefix(line, "!") {
				if strings.HasPrefix(line, "!/ADDON!") {
					break // Ignore the rest of the message.
				}
				if strings.Contains(line, "!URG!") {
					result.Urgent = true
				}
				if result.AddonName == "" {
					if match := re(`!([^!]*)!$`).FindStringSubmatch(line); match != nil {
						result.AddonName = match[1]
					}
				}
			} else if strings.HasPrefix(line, "#") {
				if match := re(`^#\s*(?:T|FORMFILENAME):(.*)`).FindStringSubmatch(line); match != nil {
					result.FormType = strings.TrimSpace(match[1])
				} else if match = re(`^#\s*(?:V|VERSION):(.*)`).FindStringSubmatch(line); match != nil {
					result.AddonVersion = strings.TrimSpace(match[1])
				}
			} else if match := re(`:\s*\[`).FindStringIndex(line); match != nil {
				fieldName = line[:match[0]]
				line = line[match[1]-1:]
			}
		}
		if fieldName != "" {
			fieldValue += line
			if value, done := unbracketData(fieldValue); done {
				result.Fields[toShortName(fieldName)] = value
				fieldName, fieldValue = "", ""
			}
		}
	}
	if messageContainsAForm(result, nil) {
		body = message
	} else {
		// It didn't contain a form, after all.  body is not repaired.
		clear(result.Fields)
	}
	if headers != "" {
		repaired = headers + eol + eol + body
	} else {
		repaired = body
	}
	return result, repaired
}

// Remove any line breaks that the BBS might have inserted into lines that look
// like form data.
func repairMessage(msg string) string {
	if msg == "" {
		return msg
	}
	// JNOS inserts a line break after every 127 bytes in a line.  It might
	// insert a break within a multi-byte UTF-8 character.  Remove any line
	// breaks that were inserted into the message.
	var message, partial string
	// Append lines into partial until a complete line is found; then append
	// partial to message.
	for _, line := range strings.Split(msg, eol) {
		// Data should have the form name: [value]
		// A ] within the value is escaped as `]
		// The value part may end with any of:
		// ] preceded by neither ` nor ]
		// `]] if the value ends with ]
		// `]]] if the value ends with `
		if re("`\\]\\]\\s*$").MatchString(partial) {
			// This is the end of a value.
			if !re(`\s$`).MatchString(partial) {
				// Look to see if there's another ] on the next line.
				if !re(`^\]\s*$`).MatchString(line) {
					// The value ended with ] represented as `]]
					message += partial + eol
					partial = ""
				} // else the value ended with ` represented as `]]]
			}
		}
		if partial == "" && // this is the first line of a group
			(!re(`:\s*\[`).MatchString(line) || // no name
				re(`^[!#]`).MatchString(line)) { // boundary marker or comment
			message += line + eol
		} else {
			partial += line
			if re("(?:`\\]\\]\\]|[^`\\]]\\])\\s*$").MatchString(partial) {
				message += partial + eol
				partial = ""
			}
		}
		// JNOS sometimes appends junk to the end of a message.
		// One observed case was "You have new messages."
		// So ignore anything after !/ADDON!.
		if partial == "" && re(`^\s*!/ADDON!\s*$`).MatchString(line) {
			break
		}
	}
	if partial != "" {
		message += partial + eol
	}
	return message
}

func unbracketData(data string) (v string, done bool) {
	if match := re(`]\s*$`).FindStringIndex(data); match != nil {
		// Remove the enclosing brackets and trailing whitespace.
		data = data[1:match[0]]
		if !strings.HasSuffix(data, "`") {
			// data is complete
			data = strings.TrimSuffix(data, "]]")
			// Un-escape the brackets with data:
			return strings.ReplaceAll(data, "`]", "]"), true
		}
	}
	return "", false
}

func toShortName(fieldName string) string {
	if idx := strings.LastIndexByte(fieldName, '.'); idx >= 0 {
		return fieldName[:idx+1]
	}
	return fieldName
}

func messageContainsAForm(message *ParsedEmail, environment map[string]string) bool {
	return (message != nil && message.FormType != "") ||
		(environment != nil && environment["ADDON_MSG_TYPE"] != "")
}

func AsciifyHeader(subject string) string {
	subject = re(`[\r\n]`).ReplaceAllLiteralString(subject, " ")
	subject = re(`[^ -~]`).ReplaceAllLiteralString(subject, "~")
	return subject
}

var reMap = make(map[string]*regexp.Regexp)

func re(s string) *regexp.Regexp {
	if re, ok := reMap[s]; ok {
		return re
	}
	re := regexp.MustCompile(s)
	reMap[s] = re
	return re
}

var eolReplaceRE = regexp.MustCompile(`^\n|([^\r])\n|\r\n?`)

func toEOL(s string) string {
	s = eolReplaceRE.ReplaceAllString(s, "$1\r\n")
	s = eolReplaceRE.ReplaceAllString(s, "$1\r\n")
	// You have to do it twice to handle multiple consecutive blank lines.
	return s
}
