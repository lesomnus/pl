package funcs

import (
	"fmt"
	"regexp"
)

type RegexMatch struct {
	Source  string
	ByIndex []string
	ByName  map[string]string
}

func (r *RegexMatch) String() string {
	return r.Source
}

func Regex(expr string, ss ...string) ([]*RegexMatch, error) {
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	rst := make([]*RegexMatch, 0, len(ss))
	for _, s := range ss {
		matched := pattern.FindStringSubmatch(s)
		if matched == nil {
			continue
		}

		match := &RegexMatch{
			Source:  s,
			ByIndex: matched[1:],
			ByName:  make(map[string]string),
		}

		for i, name := range pattern.SubexpNames()[1:] {
			if name == "" {
				continue
			}
			match.ByName[name] = matched[i+1]
		}

		rst = append(rst, match)
	}

	return rst, nil
}
