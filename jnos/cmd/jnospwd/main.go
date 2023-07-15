package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/wppsvr/config"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: jnospwd bbsname mailbox challenge\n")
		os.Exit(2)
	}
	allmsg.Register()
	config.Read()
	bbs := config.Get().BBSes[strings.ToUpper(os.Args[1])]
	if bbs == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no such BBS %q\n", os.Args[1])
		os.Exit(1)
	}
	pwd := bbs.Passwords[strings.ToUpper(os.Args[2])]
	if pwd == "" {
		fmt.Fprintf(os.Stderr, "ERROR: no password for mailbox %q on BBS %q\n", os.Args[2], os.Args[1])
		os.Exit(1)
	}
	challenge, err := strconv.ParseUint(os.Args[3], 16, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid challenge %q\n", os.Args[3])
		os.Exit(1)
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(challenge))
	buf = append(buf, pwd...)
	sum := md5.Sum(buf)
	fmt.Println(hex.EncodeToString(sum[:]))
}
