package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UAT string

func (u UAT) Version(_ ...Version) Version {
	parts := strings.Split(string(u), "-")
	if len(parts) != 2 {
		return NewUATTag(1)
	}
	if parts[0] != "uat" {
		return NewUATTag(1)
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) != 2 {
		return NewUATTag(1)
	}
	if parts[0] == dateStringUATFormat(time.Now()) {
		start, err := strconv.Atoi(parts[1])
		if err != nil {
			return NewUATTag(1)
		}
		return NewUATTag(1 + start)
	}
	return NewUATTag(1)
}

func NewUATTag(num int) UAT {
	return UAT(fmt.Sprintf("uat-%s.%d", dateStringUATFormat(time.Now()), num))
}

func dateStringUATFormat(t time.Time) string {
	return strings.ReplaceAll(t.Format(time.DateOnly), "-", "")
}
