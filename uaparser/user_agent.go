package uaparser

import "strings"

type UserAgent struct {
	Family string
	Major  string
	Minor  string
	Patch  string
}

func (parser *uaParser) Match(line string, ua *UserAgent) {
	m := parser.Reg.MatcherString(line, 0)
	if m.Matches() {
		ua.Family = strings.TrimSpace(expandString(parser.FamilyReplacement, m))
		ua.Major = strings.TrimSpace(expandString(parser.V1Replacement, m))
		ua.Minor = strings.TrimSpace(expandString(parser.V2Replacement, m))
		ua.Patch = strings.TrimSpace(expandString(parser.V3Replacement, m))
	}
}

func (ua *UserAgent) ToString() string {
	var str string
	if ua.Family != "" {
		str += ua.Family
	}
	version := ua.ToVersionString()
	if version != "" {
		str += " " + version
	}
	return str
}

func (ua *UserAgent) ToVersionString() string {
	var version string
	if ua.Major != "" {
		version += ua.Major
	}
	if ua.Minor != "" {
		version += "." + ua.Minor
	}
	if ua.Patch != "" {
		version += "." + ua.Patch
	}
	return version
}
