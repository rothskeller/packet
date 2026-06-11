package main

import (
	"encoding/json"
	"os"

	"github.com/rothskeller/packet/jnos/tnc"
)

func main() {
	t := tnc.Get(os.Args[1])
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(t)
}
