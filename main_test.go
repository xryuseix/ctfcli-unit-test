package main

import (
	"testing"
)

func AssertFlagArrayEqual(t *testing.T, expected []Flag, actual []Flag) {
	if len(expected) != len(actual) {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	for i := range expected {
		if expected[i].flag != actual[i].flag {
			t.Errorf("expected idx %v is %v, actual %v", i, expected[i], actual[i])
		}
		if expected[i].fail != actual[i].fail {
			t.Errorf("expected idx %v is %v, actual %v", i, expected[i], actual[i])
		}
	}
}

func TestRemoveComment(t *testing.T) {
	lines := []string{
		"#comment",
		"A#",
		"A#comment",
		"A #",
		"A  # ",
		"A # a",
		"A\\#escaped",
		"\\#escaped#comment\\#not-escaped",
	}
	expected := []string{
		"",
		"A",
		"A",
		"A",
		"A",
		"A",
		"A#escaped",
		"#escaped",
	}
	for i := range lines {
		actual := RemoveComment(lines[i])
		if actual != expected[i] {
			t.Errorf("expected %v, actual %v", expected[i], actual)
		}
	}
}

func TestRemoveFailFlag(t *testing.T) {
	lines := []string{
		"A",
		"!",
		"!A",
		"!!A",
		"A{great!}",
		"!A#comment",
	}
	type LineWithFail struct {
		line string
		fail bool
	}
	expected := []LineWithFail{
		{line: "A", fail: false},
		{line: "", fail: true},
		{line: "A", fail: true},
		{line: "!A", fail: true},
		{line: "A{great!}", fail: false},
		{line: "A#comment", fail: true},
	}
	for i := range lines {
		actualLine, actualFail := RemoveFailFlag(lines[i])
		if actualLine != expected[i].line || actualFail != expected[i].fail {
			t.Errorf("expected %v, actual %v", expected[i], LineWithFail{line: actualLine, fail: actualFail})
		}
	}
}

func TestParseFlag(t *testing.T) {
	content := []byte(
		"\n\nflag{flag1}\nflag{flag2}\n\nflag{flag3} # comment\nflag{flag4}# comment2\nflag{flag5}#comment3\n#comment_line\nflag{flag6\\#escaped}\nflag{flag7\\\\#escaped2}\n!flag{flag8_assert_fail}#comment\n#comment_line\n\n")
	expected := Flags{
		Flag{flag: "flag{flag1}", fail: false},
		Flag{flag: "flag{flag2}", fail: false},
		Flag{flag: "flag{flag3}", fail: false},
		Flag{flag: "flag{flag4}", fail: false},
		Flag{flag: "flag{flag5}", fail: false},
		Flag{flag: "flag{flag6#escaped}", fail: false},
		Flag{flag: "flag{flag7\\#escaped2}", fail: false},
		Flag{flag: "flag{flag8_assert_fail}", fail: true},
	}
	actual := ParseFlag(content)
	AssertFlagArrayEqual(t, expected, actual)
}

func TestParseChall(t *testing.T) {
	filePath := "example/challenge.yml"
	content := []byte(`
name: "Challenge 1"
description: "This is a challenge"
flags:
    - flag{flag1}
    - flag{flag2}
    - {
        type: "static",
        content: "flag{flag3}",
    }
    - {
        type: "regex",
        content: "flag{.*flag4.*}",
        data: "case_insensitive",
    }
`)
	expected := Challenge{
		Flags: []YamlFlag{
			{
				Type:    "static",
				Content: "flag{flag1}",
				Data:    "case_sensitive",
			},
			{
				Type:    "static",
				Content: "flag{flag2}",
				Data:    "case_sensitive",
			},
			{
				Type:    "static",
				Content: "flag{flag3}",
				Data:    "case_sensitive",
			},
			{
				Type:    "regex",
				Content: "flag{.*flag4.*}",
				Data:    "case_insensitive",
			},
		},
	}
	actual, err := ParseChall(filePath, content)
	if err != nil {
		t.Errorf("error parsing challenge: %v", err)
	}
	if len(expected.Flags) != len(actual.Flags) {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	for i := range expected.Flags {
		if expected.Flags[i] != actual.Flags[i] {
			t.Errorf("expected idx %v is %v, actual %v", i, expected.Flags[i], actual.Flags[i])
		}
	}
}

func TestUnitTestOk(t *testing.T) {
	challs := map[string]Challenge{
		"challenge1": {
			Flags: []YamlFlag{
				{
					Type:    "regex",
					Content: "flag{.*flag1.*}",
					Data:    "case_insensitive",
				},
			},
		},
	}
	flags := map[string]Flags{
		"challenge1": {
			Flag{
				flag: "flag{flag1}",
				fail: false,
			},
		},
	}
	isErr := UnitTest(challs, flags)
	if isErr {
		t.Errorf("expected false, actual true")
	}
}

func TestUnitTestNg(t *testing.T) {
	challs := map[string]Challenge{
		"challenge1": {
			Flags: []YamlFlag{
				{
					Type:    "regex",
					Content: "flag[123]",
					Data:    "case_insensitive",
				},
			},
		},
	}
	flags := map[string]Flags{
		"challenge1": {Flag{flag: "flag{flag1}", fail: false}},
	}
	isErr := UnitTest(challs, flags)
	if !isErr {
		t.Errorf("expected error, actual no error")
	}
}
