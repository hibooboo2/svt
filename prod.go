package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PROD string

func (p PROD) Version(_ ...Version) Version {
	parts := []string{string(p[0]), string(p[1:])}
	if len(parts) != 2 {
		return NewProdTag(1)
	}
	if parts[0] != "r" {
		return NewProdTag(1)
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) != 2 {
		return NewProdTag(1)
	}
	if parts[0] == dateStringUATFormat(time.Now()) {
		start, err := strconv.Atoi(parts[1])
		if err != nil {
			return NewProdTag(1)
		}
		return NewProdTag(1 + start)
	}
	return NewProdTag(1)
}

func NewProdTag(num int) PROD {
	return PROD(fmt.Sprintf("r%s.%d", dateStringUATFormat(time.Now()), num))
}
