package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lastTag := flag.Bool("last-tag", false, "the last tag")
	mode := flag.String("mode", "dev", "mode is what kind of version system you want to increment valid choices are ['dev','uat','prod','img']")
	flag.Parse()

	args := flag.Args()

	if *lastTag {
		lines := bufio.NewReader(os.Stdin)
		//list out the tags and their versions
		// convert to ssemver based on type
		// get the latest version
		// print the latest version
		versions := []Version{}
		for {
			line, err := lines.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			var v Version
			switch *mode {
			case "dev":
				v = SemVer(line)
			case "uat":
				v = UAT(line)
			case "prod":
				v = PROD(line)
			case "img":
				v = SemVer(line[len("img-"):])
			default:
				panic("Invalid mode")
			}
			versions = append(versions, v)
		}

		return
	}

	switch len(args) {
	case 1, 2:
	default:
		panic(args)
	}

	var v Version
	var v2 Version
	switch *mode {
	case "dev":
		if len(os.Args) < 2 {
			panic(args)
		}
		v = SemVer(args[0])
		v2 = SemVer(args[1])
	case "uat":
		v = UAT(args[0])
	case "prod", "img":
		v = PROD(args[0])
	}
	v = v.Version(v2)
	switch *mode {
	case "img":
		imgTag := fmt.Sprintf("%s", v)
		fmt.Printf("img-%s", imgTag[1:])
	default:
		fmt.Println(v)
	}
}

type Version interface {
	Version(...Version) Version
}
