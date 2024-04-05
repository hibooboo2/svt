package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	version := os.Args[1]
	v := SemVer(version)
	v = v.Add(SemVer(os.Args[2]))
	fmt.Println(v)
}

type SemVer string

func (sv SemVer) Add(verToAdd SemVer) SemVer {
	return MakeNewSemVer(sv.MajorAdd(verToAdd.Major()), sv.MinorAdd(verToAdd.Minor()), sv.PatchAdd(verToAdd.Patch()))
}

func (sv SemVer) MajorAdd(toAdd int) int {
	return toAdd + sv.Major()
}

func (sv SemVer) Major() int {
	parts := strings.Split(string(sv), ".")
	if len(parts) != 3 {
		panic("Invalid semVer")
	}
	version := strings.TrimPrefix(parts[0], "v")
	v, err := strconv.Atoi(version)
	if err != nil {
		panic(fmt.Errorf("failed to parse number from %q: %w", version, err))
	}
	return v
}

func (sv SemVer) MinorAdd(toAdd int) int {
	return toAdd + sv.Minor()
}

func (sv SemVer) Minor() int {
	parts := strings.Split(string(sv), ".")
	if len(parts) != 3 {
		panic("Invalid semVer")
	}
	v, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Errorf("failed to parse number from %q: %w", err))
	}
	return v
}

func (sv SemVer) PatchAdd(toAdd int) int {
	return toAdd + sv.Patch()
}

func (sv SemVer) Patch() int {
	parts := strings.Split(string(sv), ".")
	if len(parts) != 3 {
		panic("Invalid semVer")
	}
	v, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(fmt.Errorf("failed to parse number from %q: %w", err))
	}
	return v
}

func MakeNewSemVer(maj, min, patch int) SemVer {
	return SemVer(fmt.Sprintf("v%d.%d.%d", maj, min, patch))
}
