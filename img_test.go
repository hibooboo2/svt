package main

import "testing"

func TestImg(t *testing.T) {
	tag := PROD("img-20260715.2")

	var v2 Version

	q := tag.Version(v2)

	_ = q
}

func TestProd(t *testing.T) {
	tag := PROD("r20260715.2")

	var v2 Version

	q := tag.Version(v2)

	_ = q
}
