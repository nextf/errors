package errors

import "regexp"

type ConstError string

var regForErrCode *regexp.Regexp = regexp.MustCompile(`^\s*\[([A-Za-z0-9_-]+)\]\s*(.*)\s*$`)

func (e ConstError) Error() string {
	if group := regForErrCode.FindStringSubmatch(string(e)); len(group) >= 3 {
		return group[2]
	}
	return string(e)
}

func (e ConstError) Code() string {
	if group := regForErrCode.FindStringSubmatch(string(e)); len(group) >= 2 {
		return group[1]
	}
	return ""
}

func (e ConstError) Match(key interface{}) bool {
	if key == nil {
		return false
	}
	switch x := key.(type) {
	case string:
		return e.Code() == x
	case interface{ MatchString(s string) bool }:
		return x.MatchString(e.Code())
	}
	return false
}
