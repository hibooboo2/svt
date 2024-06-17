package main

import (
	"flag"
	"fmt"
)

func main() {
	mode := flag.String("mode", "dev", "mode is what kind of version system you want to increment valid choices are ['dev','uat']")
	flag.Parse()

	args := flag.Args()

	switch len(args) {
	case 1, 2:
	default:
		panic(args)
	}

	var v Version
	var v2 Version
	switch *mode {
	case "dev":
		v = SemVer(args[0])
		v2 = SemVer(args[1])
	case "uat":
		v = UAT(args[0])
	}
	v = v.Version(v2)
	fmt.Println(v)
}

type Version interface {
	Version(...Version) Version
}
