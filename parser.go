package main

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorReset  = "\x1b[0m"
)

var (
	reComment  = regexp.MustCompile(`\s*#[^#\s]*`)
	reFailFlag = regexp.MustCompile(`^\s*!`)
	reNewlines = regexp.MustCompile(`\n+`)
)

type YamlFlag struct {
	Type    string `yaml:"type,omitempty"`
	Content string `yaml:"content,omitempty"`
	Data    string `yaml:"data,omitempty"`
}

type ChallYaml struct {
	Type  string        `yaml:"type,omitempty"`
	Flags []any `yaml:"flags"`
}

type Challenge struct {
	Type  string
	Flags []YamlFlag
}

type Config struct {
	Genre []string `yaml:"genre,omitempty"`
}

func ParseChall(filePath string, content []byte) (Challenge, error) {
	text := string(content)
	if strings.Contains(text, "\t") {
		fmt.Printf("%sWarning%s: TAB is not recommended in the YAML file (%v).\nPlease see the FAQ: https://yaml.org/faq.html\n", colorYellow, colorReset, filePath)
		content = []byte(strings.ReplaceAll(text, "\t", "    "))
	}

	chall := Challenge{
		Type:  "",
		Flags: []YamlFlag{},
	}

	var challYaml ChallYaml
	err := yaml.Unmarshal(content, &challYaml)
	if err != nil {
		fmt.Printf("%sError%s: unmarshalling the file %v: %v\n", colorRed, colorReset, filePath, err)
		return chall, err
	}

	chall.Type = challYaml.Type

	for _, flag := range challYaml.Flags {
		switch flag := flag.(type) {
		case string:
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    "static",
				Content: flag,
				Data:    "case_sensitive",
			})
		case map[string]any:
			flagType, ok := flag["type"]
			if !ok {
				flagType = "static"
			}
			flagContent, ok := flag["content"]
			if !ok {
				return chall, fmt.Errorf("%sError%s: flag content is not specified", colorRed, colorReset)
			}
			flagData, ok := flag["data"]
			if !ok {
				flagData = "case_sensitive"
			}
			ft, ok := flagType.(string)
			if !ok {
				return chall, fmt.Errorf("%sError%s: invalid flag type: %v", colorRed, colorReset, flagType)
			}
			fc, ok := flagContent.(string)
			if !ok {
				return chall, fmt.Errorf("%sError%s: invalid flag content: %v", colorRed, colorReset, flagContent)
			}
			fd, ok := flagData.(string)
			if !ok {
				return chall, fmt.Errorf("%sError%s: invalid flag data: %v", colorRed, colorReset, flagData)
			}
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    ft,
				Content: fc,
				Data:    fd,
			})
		default:
			return chall, fmt.Errorf("%sError%s: unknown flag type: %v", colorRed, colorReset, flag)
		}
	}

	return chall, nil
}

type Flag struct {
	flag string
	fail bool // if true, the flag-check must be failed
}

type Flags = []Flag

func RemoveComment(line string) string {
	matches := reComment.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		return line
	}
	escaped := 0
	for _, match := range matches {
		if match[0] > 0 && line[match[0]-1] == '\\' {
			line = line[:match[0]-1-escaped] + line[match[0]-escaped:]
			escaped++
			continue
		}
		line = line[:match[0]-escaped]
		break
	}
	return line
}

func RemoveFailFlag(line string) (string, bool) {
	idx := reFailFlag.FindStringIndex(line)
	if idx == nil {
		return line, false
	}
	return line[idx[1]:], true
}

func ParseFlag(content []byte) Flags {
	text := reNewlines.ReplaceAllString(string(content), "\n")
	text = strings.Trim(text, "\n")
	var flags Flags
	for _, flag := range strings.Split(text, "\n") {
		flag = RemoveComment(flag)
		if flag == "" {
			continue
		}
		flag, fail := RemoveFailFlag(flag)
		flags = append(flags, Flag{flag: flag, fail: fail})
	}
	return flags
}
