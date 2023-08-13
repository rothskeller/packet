package envelope

// This file contains code for parsing addresses and address lists.

import (
	"errors"
	"strings"
)

// Address is an address that can appear in the header of a message.  It
// consists of an optional name comment followed by an xxx@yyy address.
type Address struct {
	Name    string
	Address string
}

func (addr *Address) String() string {
	var local, domain, _ = strings.Cut(addr.Address, "@")
	// The local part may need to be quoted.
	for i, r := range local {
		if isAText(r) {
			continue
		}
		if r == '.' && i > 0 && local[i-1] != '.' && i < len(local)-1 {
			continue
		}
		local = quoteString(local)
		break
	}
	if domain != "" {
		local += "@" + domain
	}
	if addr.Name == "" {
		return local
	}
	// The name may also need to be quoted.
	var name = addr.Name
	for _, r := range name {
		if isAText(r) || isWhitespace(r) {
			continue
		}
		name = quoteString(name)
		break
	}
	return name + " <" + local + ">"
}

func quoteString(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, r := range s {
		if !isQText(r) && !isWhitespace(r) {
			sb.WriteByte('\\')
		}
		sb.WriteRune(r)
	}
	sb.WriteByte('"')
	return sb.String()
}

/*
This is the grammar we are parsing.  It is exactly as described in RFC-5322
except for one enhancement
  (1) addr-spec can have a local-part without a domain
and two limitations
  (1) no group list syntax
  (2) no obsolete syntax options
The enhancement allows us to send packet messages to other mailboxes on the
same BBS without an @bbs suffix.  This style of addressing is discouraged but
common.

   address-list    =   (address *("," address))

   address         =   mailbox

   mailbox         =   name-addr / addr-spec

   name-addr       =   [display-name] angle-addr

   angle-addr      =   [CFWS] "<" addr-spec ">" [CFWS]

   display-name    =   phrase

   addr-spec       =   local-part "@" domain /
                       local-part

   local-part      =   dot-atom / quoted-string

   domain          =   dot-atom / domain-literal

   domain-literal  =   [CFWS] "[" *([FWS] dtext) [FWS] "]" [CFWS]

   dtext           =   %d33-90 /          ; Printable US-ASCII
                       %d94-126           ;  characters not including
                                          ;  "[", "]", or "\"

   FWS             =   ([*WSP CRLF] 1*WSP)
                                          ; Folding white space

   ctext           =   %d33-39 /          ; Printable US-ASCII
                       %d42-91 /          ;  characters not including
                       %d93-126           ;  "(", ")", or "\"

   ccontent        =   ctext / quoted-pair / comment

   quoted-pair     =   ("\" (VCHAR / WSP))

   comment         =   "(" *([FWS] ccontent) [FWS] ")"

   CFWS            =   (1*([FWS] comment) [FWS]) / FWS

   atext           =   ALPHA / DIGIT /    ; Printable US-ASCII
                       "!" / "#" /        ;  characters not including
                       "$" / "%" /        ;  specials.  Used for atoms.
                       "&" / "'" /
                       "*" / "+" /
                       "-" / "/" /
                       "=" / "?" /
                       "^" / "_" /
                       "`" / "{" /
                       "|" / "}" /
                       "~"

   atom            =   [CFWS] 1*atext [CFWS]

   dot-atom-text   =   1*atext *("." 1*atext)

   dot-atom        =   [CFWS] dot-atom-text [CFWS]

   qtext           =   %d33 /             ; Printable US-ASCII
                       %d35-91 /          ;  characters not including
                       %d93-126           ;  "\" or the quote character

   qcontent        =   qtext / quoted-pair

   quoted-string   =   [CFWS]
                       DQUOTE *([FWS] qcontent) [FWS] DQUOTE
                       [CFWS]

   phrase          =   1*word

   word            =   atom / quoted-string

*/

func ParseAddressList(s string) (addrs []*Address, err error) {
	_, s = parseWhitespace(s)
	if s == "" {
		return nil, nil
	}
	if addr, rest, ok := parseAddress(s); !ok {
		return nil, errors.New("invalid address list")
	} else {
		addrs = append(addrs, addr)
		s = rest
	}
	for s != "" {
		if s[0] != ',' {
			return nil, errors.New("invalid address list")
		}
		s = s[1:]
		if addr, rest, ok := parseAddress(s); !ok {
			return nil, errors.New("invalid address list")
		} else {
			addrs = append(addrs, addr)
			s = rest
		}
	}
	return addrs, nil
}

func ParseAddress(s string) (*Address, error) {
	addr, rest, ok := parseAddress(s)
	if !ok || rest != "" {
		return nil, errors.New("invalid address")
	}
	return addr, nil
}

func parseAddress(s string) (addr *Address, rest string, ok bool) {
	if addr, rest, ok = parseNameAddr(s); ok {
		return
	}
	return parseAddrSpec(s)
}

func parseNameAddr(s string) (addr *Address, rest string, ok bool) {
	var a Address

	for {
		word, rest, ok := parseWord(s)
		if !ok {
			break
		}
		if a.Name != "" {
			a.Name += " "
		}
		a.Name += word
		s = rest
	}
	if s == "" || s[0] != '<' {
		return nil, "", false
	}
	s = s[1:]
	if addr, rest, ok := parseAddrSpec(s); !ok {
		return nil, "", false
	} else {
		a.Address = addr.Address
		s = rest
	}
	if s == "" || s[0] != '>' {
		return nil, "", false
	}
	s = s[1:]
	return &a, s, true
}

func parseAddrSpec(s string) (addr *Address, rest string, ok bool) {
	var a Address

	if lp, rest, ok := parseLocalPart(s); ok {
		a.Address = lp
		s = rest
	} else {
		return nil, "", false
	}
	if s == "" || s[0] != '@' {
		return &a, s, true
	}
	if dom, rest, ok := parseDomain(s[1:]); ok {
		a.Address += "@" + dom
		s = rest
		return &a, rest, true
	}
	return nil, "", false
}

func parseLocalPart(s string) (lp, rest string, ok bool) {
	if da, rest, ok := parseDotAtom(s); ok {
		return da, rest, ok
	}
	return parseQuotedString(s)
}

func parseDomain(s string) (dom, rest string, ok bool) {
	if da, rest, ok := parseDotAtom(s); ok {
		return da, rest, ok
	}
	return parseDomainLiteral(s)
}

func parseDotAtom(s string) (da, rest string, ok bool) {
	_, s = parseCommentsWhitespace(s)
	if atextrun, rest, ok := parseATextRun(s); ok {
		da = atextrun
		s = rest
	} else {
		return "", "", false
	}
	for s != "" && s[0] == '.' {
		if atextrun, rest, ok := parseATextRun(s[1:]); ok {
			da += "." + atextrun
			s = rest
		} else {
			break
		}
	}
	return da, s, true
}

func parseATextRun(s string) (run, rest string, ok bool) {
	idx := strings.IndexFunc(s, func(r rune) bool {
		return r < 256 && strings.IndexByte(atextchars, byte(r)) < 0
	})
	if idx == 0 || s == "" {
		return "", "", false
	}
	if idx < 0 {
		return s, "", true
	}
	return s[:idx], s[idx:], true
}

func parseDomainLiteral(s string) (dl, rest string, ok bool) {
	_, s = parseCommentsWhitespace(s)
	if s == "" || s[0] != '[' {
		return "", "", false
	}
	idx := strings.IndexFunc(s, func(r rune) bool {
		return r != 9 && (r < 32 || (r > 90 && r < 94) || r > 126)
	})
	if idx < 0 || s[idx] != ']' {
		return "", "", false
	}
	dl = s[:idx+1]
	_, s = parseCommentsWhitespace(s[idx+1:])
	return dl, s, true
}

func parseWord(s string) (word, rest string, ok bool) {
	if word, rest, ok = parseAtom(s); ok {
		return
	}
	return parseQuotedString(s)
}

const atextchars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-/=?^_`{|}~"

func parseAtom(s string) (atom, rest string, ok bool) {
	_, s = parseCommentsWhitespace(s)
	if s == "" || strings.IndexByte(atextchars, s[0]) < 0 {
		return "", "", false
	}
	idx := strings.IndexFunc(s, func(r rune) bool {
		return strings.IndexByte(atextchars, byte(r)) < 0
	})
	if idx < 0 {
		atom, s = s, ""
	} else {
		atom, s = s[:idx], s[idx:]
	}
	_, s = parseCommentsWhitespace(s)
	return atom, s, true
}

func parseQuotedString(s string) (qs, rest string, ok bool) {
	_, s = parseCommentsWhitespace(s)
	if s == "" || s[0] != '"' {
		return "", "", false
	}
	qs, s = s[:1], s[1:]
	for {
		if sp, rest := parseWhitespace(s); sp != "" {
			qs += sp
			s = rest
		}
		if qc, rest, ok := parseQContent(s); ok {
			qs += qc
			s = rest
		} else {
			break
		}
	}
	if s == "" || s[0] != '"' {
		return "", "", false
	}
	qs, s = qs+s[:1], s[1:]
	return qs, s, true
}

func parseQContent(s string) (qc, rest string, ok bool) {
	for {
		if qtrun, rest, ok := parseQTextRun(s); ok {
			qc += qtrun
			s = rest
		} else if qp, rest, ok := parseQuotedPair(s); ok {
			qc += qp
			s = rest
		} else {
			break
		}
	}
	if qc == "" {
		return "", "", false
	}
	return qc, s, true
}

func parseQTextRun(s string) (qtrun, rest string, ok bool) {
	idx := strings.IndexFunc(s, func(r rune) bool {
		return r != 33 && (r < 35 || r > 91) && (r < 93 || r > 126)
	})
	if idx == 0 {
		return "", "", false
	}
	if idx < 0 {
		return s, "", true
	}
	return s[:idx], s[idx:], true
}

func parseCommentsWhitespace(s string) (com, rest string) {
	for {
		if ws, rest := parseWhitespace(s); ws != "" {
			com += ws
			s = rest
		}
		if c, rest, ok := parseComment(s); ok {
			com += c
			s = rest
		} else {
			break
		}
	}
	return com, s
}

func parseComment(s string) (com, rest string, ok bool) {
	if s == "" || s[0] != '(' {
		return "", "", false
	}
	s = s[1:]
	for {
		var ws string
		ws, s = parseWhitespace(s)
		com += ws
		if cc, rest, ok := parseCContent(s); ok {
			com += cc
			s = rest
		} else {
			break
		}
	}
	if s == "" || s[0] != ')' {
		return "", "", false
	}
	return com, s[1:], true
}

func parseCContent(s string) (cc, rest string, ok bool) {
	if ct, rest, ok := parseCTextRun(s); ok {
		return ct, rest, true
	} else if qp, rest, ok := parseQuotedPair(s); ok {
		return qp, rest, true
	} else if com, rest, ok := parseComment(s); ok {
		return com, rest, true
	}
	return "", "", false
}

func parseCTextRun(s string) (ctrun, rest string, ok bool) {
	idx := strings.IndexFunc(s, func(r rune) bool {
		return (r < 33 || r > 39) && (r < 42 || r > 91) && (r < 93 || r > 126)
	})
	if idx == 0 {
		return "", "", false
	}
	if idx < 0 {
		return s, "", true
	}
	return s[:idx], s[idx:], true
}

func parseQuotedPair(s string) (pair, rest string, ok bool) {
	if len(s) < 2 || s[0] != '\\' {
		return "", "", false
	}
	if s[1] == 9 || (s[1] >= 32 && s[1] <= 126) {
		return s[:2], s[2:], true
	}
	return "", "", false
}

func parseWhitespace(s string) (ws, rest string) {
	idx := strings.IndexFunc(s, func(r rune) bool { return r != ' ' && r != '\t' })
	if idx < 0 {
		return "", rest
	}
	return s[:idx], s[idx:]
}

func isAText(r rune) bool {
	switch r {
	case '.', '(', ')', '[', ']', ';', '@', '\\', ',', '<', '>', '"', ':':
		return false
	}
	return isVChar(r)
}

func isQText(r rune) bool {
	return r != '"' && r != '\\' && isVChar(r)
}

func isVChar(r rune) bool {
	return '!' <= r && r <= '~'
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}
