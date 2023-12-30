package main

import (
	"testing"
)

func AssertArrayEqual(t *testing.T, expected []string, actual []string) {
	if len(expected) != len(actual) {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("expected idx %v is %v, actual %v", i, expected[i], actual[i])
		}
	}
}

func TestParseFlag(t *testing.T) {
	filePath := "example/flag.txt"
	content := []byte("\n\nflag{flag1}\nflag{flag2}\n\nflag{flag3}\n\n")
	expected := []string{"flag{flag1}", "flag{flag2}", "flag{flag3}"}
	actual := ParseFlag(filePath, content)
	AssertArrayEqual(t, expected, actual)
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
