package main

import (
	"flag"
)

// Flags keeps settings from commandline
type Flags struct {
	AutoLink bool
	Thomas   bool
	Watch    bool
}

var flags Flags

func init() {
	flag.BoolVar(&flags.AutoLink, "a", false, "autolink")
	flag.BoolVar(&flags.Thomas, "t", false, "thomas")
	flag.BoolVar(&flags.Watch, "w", false, "watch")
	flag.Parse()
}
