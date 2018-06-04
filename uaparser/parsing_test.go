package uaparser

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var testParser *Parser

func init() {
	var err error
	testParser, err = New("../uap-core/regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

func TestOSParsing(t *testing.T) {
	var YAMLTestCases struct {
		Cases []struct {
			Input    string `yaml:"user_agent_string"`
			Expected Os     `yaml:",inline"`
		} `yaml:"test_cases"`
	}

	for _, filepath := range []string{"../uap-core/test_resources/additional_os_tests.yaml", "../uap-core/tests/test_os.yaml"} {
		unmarshalResourceTestFile(filepath, &YAMLTestCases)

		for _, c := range YAMLTestCases.Cases {
			actual := testParser.ParseOs(c.Input)
			if got, want := *actual, c.Expected; got != want {
				t.Fatalf("got %#v\n want %#v\n", got, want)
			}
		}
	}
}

func TestReadsInteralYAML(t *testing.T) {
	_ = NewFromSaved() // should not panic
}

func TestUAParsing(t *testing.T) {
	t.Parallel()

	var YAMLTestCases struct {
		Cases []struct {
			Input    string    `yaml:"user_agent_string"`
			Expected UserAgent `yaml:",inline"`
		} `yaml:"test_cases"`
	}

	for _, filepath := range []string{"../uap-core/test_resources/pgts_browser_list.yaml", "../uap-core/test_resources/firefox_user_agent_strings.yaml", "../uap-core/tests/test_ua.yaml"} {
		unmarshalResourceTestFile(filepath, &YAMLTestCases)

		for _, c := range YAMLTestCases.Cases {
			actual := testParser.ParseUserAgent(c.Input)
			if got, want := *actual, c.Expected; got != want {
				t.Fatalf("got %#v\n want %#v\n", got, want)
			}
		}
	}
}

func TestDeviceParsing(t *testing.T) {
	t.Parallel()

	var YAMLTestCases struct {
		Cases []struct {
			Input    string `yaml:"user_agent_string"`
			Expected Device `yaml:",inline"`
		} `yaml:"test_cases"`
	}

	for _, filepath := range []string{"../uap-core/tests/test_device.yaml"} {
		unmarshalResourceTestFile(filepath, &YAMLTestCases)

		for _, c := range YAMLTestCases.Cases {
			actual := testParser.ParseDevice(c.Input)
			if got, want := *actual, c.Expected; got != want {
				t.Fatalf("ua: %s, got %#v\n want %#v\n", c.Input, got, want)
			}
		}
	}
}

func TestParser_TWM(t *testing.T) {

	testSavedParser := NewFromSaved()

	tests := []struct {
		name  string
		input string
		want  *Device
	}{
		// TODO: Add test cases.
		{
			name:  "TWM A8",
			input: "Mozilla/5.0 (Linux; Android 4.4; Amazing A8 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.133 Mobile Safari/537.36",
			want: &Device{
				Family: "TWM Amazing A8",
				Brand:  "TWM",
				Model:  "Amazing A8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testSavedParser.ParseDevice(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ua: %s, got %#v\n want %#v\n", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenericParseMethodConcurrency(t *testing.T) { // go test -race -run=Concurrency
	var YAMLTestCases struct {
		Cases []struct {
			Input    string `yaml:"user_agent_string"`
			Expected Os     `yaml:",inline"`
		} `yaml:"test_cases"`
	}

	for _, filepath := range []string{"../uap-core/tests/test_os.yaml"} {
		unmarshalResourceTestFile(filepath, &YAMLTestCases)

		for _, c := range YAMLTestCases.Cases {
			actual := testParser.Parse(c.Input).Os
			if got, want := *actual, c.Expected; got != want {
				t.Fatalf("got %#v\n want %#v\n", got, want)
			}
		}
	}
}

func unmarshalResourceTestFile(filepath string, v interface{}) {
	file, err := ioutil.ReadFile(filepath)
	if nil != err {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(file, v); err != nil {
		log.Fatal(err)
	}
}
