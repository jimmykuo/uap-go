package uaparser

type Os struct {
	Family     string
	Major      string
	Minor      string
	Patch      string
	PatchMinor string `yaml:"patch_minor"`
}

func (parser *osParser) Match(line string, os *Os) {
	m := parser.Reg.MatcherString(line, 0)
	if m.Matches() {
		os.Family = expandString(parser.OSReplacement, m)
		os.Major = expandString(parser.V1Replacement, m)
		os.Minor = expandString(parser.V2Replacement, m)
		os.Patch = expandString(parser.V3Replacement, m)
		os.PatchMinor = expandString(parser.V4Replacement, m)
	}
}

func (os *Os) ToString() string {
	var str string
	if os.Family != "" {
		str += os.Family
	}
	version := os.ToVersionString()
	if version != "" {
		str += " " + version
	}
	return str
}

func (os *Os) ToVersionString() string {
	var version string
	if os.Major != "" {
		version += os.Major
	}
	if os.Minor != "" {
		version += "." + os.Minor
	}
	if os.Patch != "" {
		version += "." + os.Patch
	}
	if os.PatchMinor != "" {
		version += "." + os.PatchMinor
	}
	return version
}
