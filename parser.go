package main

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type YamlFlag struct {
	Type    string `yaml:"type,omitempty"`
	Content string `yaml:"content,omitempty"`
	Data    string `yaml:"data,omitempty"`
}

type ChallYaml struct {
	Flags []interface{} `yaml:"flags"`
}

type Challenge struct {
	Flags []YamlFlag
}

type Config struct {
	Genre []string `yaml:"genre,omitempty"`
}

func ParseChall(filePath string, content []byte) (Challenge, error) {
	text := string(content)
	if strings.Contains(text, "\t") {
		fmt.Printf("\x1b[33mWarning\x1b[0m: TAB is not recommended in the YAML file (%v).\nPlease see the FAQ: https://yaml.org/faq.html\n", filePath)
		content = []byte(strings.ReplaceAll(text, "\t", "    "))
	}

	chall := Challenge{
		Flags: []YamlFlag{},
	}

	var challYaml ChallYaml
	err := yaml.Unmarshal(content, &challYaml)
	if err != nil {
		fmt.Printf("\x1b[31mError\x1b[0m: unmarshalling the file %v: %v\n", filePath, err)
		return chall, err
	}

	for _, flag := range challYaml.Flags {
		switch flag := flag.(type) {
		case string:
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    "static",
				Content: flag,
				Data:    "case_sensitive",
			})
		case map[string]interface{}:
			flagType, ok := flag["type"]
			if !ok {
				flagType = "static"
			}
			flagContent, ok := flag["content"]
			if !ok {
				return chall, fmt.Errorf("\x1b[31mError\x1b[0m: flag content is not specified")
			}
			flagData, ok := flag["data"]
			if !ok {
				flagData = "case_sensitive"
			}
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    flagType.(string),
				Content: flagContent.(string),
				Data:    flagData.(string),
			})
		default:
			return chall, fmt.Errorf("\x1b[31mError\x1b[0m: unknown flag type: %v", flag)
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
	reg := regexp.MustCompile(`\s*#[^#]*`)
	matches := reg.FindAllStringIndex(line, -1)
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

func RemoveFailFlag(flag string) (string, bool) {
	// reg := regexp.MustCompile(`^!`)
	// return reg.ReplaceString(flag, ""), reg.MatchString(flag)
	return flag, false
}

func ParseFlag(content []byte) Flags {
	reg := regexp.MustCompile(`\n+`)
	text := reg.ReplaceAllString(string(content), "\n")
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
