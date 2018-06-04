package uaparser

import (
	"strings"
)

type Device struct {
	Family string
	Brand  string
	Model  string
}

func (parser *deviceParser) Match(line string, dvc *Device) {
	m := parser.Reg.MatcherString(line, 0)
	if m.Matches() {
		dvc.Family = strings.TrimSpace(expandString(parser.DeviceReplacement, m))
		dvc.Brand = strings.TrimSpace(expandString(parser.BrandReplacement, m))
		dvc.Model = strings.TrimSpace(expandString(parser.ModelReplacement, m))
	}
}

func (dvc *Device) ToString() string {
	return dvc.Family
}
