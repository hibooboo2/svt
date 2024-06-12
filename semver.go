package main

import (
	"fmt"
	"strconv"
	"strings"
)

type SemVer string

func MustSemVer(val SemVer) []string {
	parts := strings.Split(string(val), ".")
	if len(parts) != 3 {
		panic("Invalid semVer")
	}
	return parts
}

func (sv SemVer) Version(v ...Version) Version {
	if len(v) == 0 {
		v = append(v, SemVer("0.0.1"))
	}
	verToAdd, _ := v[0].(SemVer)
	return MakeNewSemVer(sv.MajorAdd(verToAdd.Major()), sv.MinorAdd(verToAdd.Minor()), sv.PatchAdd(verToAdd.Patch()))
}

func (sv SemVer) MajorAdd(toAdd int) int {
	return toAdd + sv.Major()
}

func (sv SemVer) Major() int {
	parts := MustSemVer(sv)
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
	parts := MustSemVer(sv)
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
	parts := MustSemVer(sv)
	v, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(fmt.Errorf("failed to parse number from %q: %w", err))
	}
	return v
}

func MakeNewSemVer(maj, min, patch int) SemVer {
	return SemVer(fmt.Sprintf("v%d.%d.%d", maj, min, patch))
}
