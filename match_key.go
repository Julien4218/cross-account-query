package main

import (
	"fmt"
	"strconv"
)

type MatchKey struct {
	Match   string
	KeyName string
}

func NewMatchKey(match string, keyName string) *MatchKey {
	return &MatchKey{Match: match, KeyName: keyName}
}

func (mk *MatchKey) Replace(value string) string {
	if mk.Match == "add1UnixDay::" {
		int64, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			// 86400000 milliseconds in 1 day
			result := int64 + 60*60*24*1000
			return fmt.Sprintf("%d", result)
		}
	}
	return value
}
