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

func TestParser_Custom(t *testing.T) {

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
		}, {
			name:  "InFocus",
			input: "Mozilla/5.0 (Linux; Android 4.4; InFocus M2_3G Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36",
			want: &Device{
				Family: "InFocus M2_3G",
				Brand:  "InFocus",
				Model:  "M2_3G",
			},
		}, {
			name:  "MIX",
			input: "Mozilla/5.0 (Linux; U; Android 6.0.1; zh-cn; MIX Build/MXB48T) AppleWebKit/537.36 (KHTML, like Gecko)Version/4.0 Chrome/37.0.0.0 MQQBrowser/7.9 Mobile Safari/537.36",
			want: &Device{
				Family: "XiaoMi",
				Brand:  "XiaoMi",
				Model:  "MIX",
			},
		}, {
			name:  "MI",
			input: "Mozilla/5.0 (Linux; U; Android 6.0.1; zh-cn; MI 4LTE Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko)Version/4.0 Chrome/37.0.0.0 MQQBrowser/7.2 Mobile Safari/537.36",
			want: &Device{
				Family: "XiaoMi MI 4LTE",
				Brand:  "XiaoMi",
				Model:  "MI 4LTE",
			},
		}, {
			name:  "XiaoMi",
			input: "Mozilla/5.0 (Linux; U; Android 5.1.1; zh-cn; MI 4S Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.146 Mobile Safari/537.36 XiaoMi/MiuiBrowser/9.1.3",
			want: &Device{
				Family: "XiaoMi MI 4S",
				Brand:  "XiaoMi",
				Model:  "MI 4S",
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
