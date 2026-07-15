package main

import "testing"

func TestImg(t *testing.T) {
	tag := PROD("img-20260715.1")

	var v2 Version

	q := tag.Version(v2)

	_ = q
}
