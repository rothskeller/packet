package main

import (
	"github.com/rothskeller/packet/message/allmsg"
	"github.com/rothskeller/packet/shell"
)

func main() {
	allmsg.Register()
	shell.Main()
}
