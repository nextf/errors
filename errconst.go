// Copyright 2022 Zhiwen<zhiwen.t@outlook.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"regexp"
	"strings"
)

type ConstError string

var regForErrCode *regexp.Regexp = regexp.MustCompile(`^\s*\[([A-Za-z0-9_-]+)\]\s*(.*)$`)

func (e ConstError) Error() string {
	if group := regForErrCode.FindStringSubmatch(string(e)); len(group) >= 3 {
		return strings.TrimSpace(group[2])
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
