package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
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

func ParseChall(filePath string, content []byte) (Challenge, error) {
	text := string(content)
	if strings.Contains(text, "\t") {
		fmt.Printf("Warning: TAB is not recommended in the YAML file (%v).\nPlease see the FAQ: https://yaml.org/faq.html\n", filePath)
		content = []byte(strings.ReplaceAll(text, "\t", "    "))
	}

	chall := Challenge{
		Flags: []YamlFlag{},
	}

	var challYaml ChallYaml
	err := yaml.Unmarshal(content, &challYaml)
	if err != nil {
		fmt.Printf("Error unmarshalling the file %v: %v\n", filePath, err)
		return chall, err
	}

	for _, flag := range challYaml.Flags {
		switch flag.(type) {
		case string:
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    "static",
				Content: flag.(string),
				Data:    "case_sensitive",
			})
		case map[string]interface{}:
			flagMap := flag.(map[string]interface{})
			flagType, ok := flagMap["type"]
			if !ok {
				flagType = "static"
			}
			flagContent, ok := flagMap["content"]
			if !ok {
				return chall, fmt.Errorf("flag content is not specified")
			}
			flagData, ok := flagMap["data"]
			if !ok {
				flagData = "case_sensitive"
			}
			chall.Flags = append(chall.Flags, YamlFlag{
				Type:    flagType.(string),
				Content: flagContent.(string),
				Data:    flagData.(string),
			})
		default:
			return chall, fmt.Errorf("unknown flag type: %v", flag)
		}
	}

	return chall, nil
}

type Flag = []string

func ParseFlag(filePath string, content []byte) Flag {
	reg := regexp.MustCompile(`\n+`)
	text := reg.ReplaceAllString(string(content), "\n")
	if strings.HasPrefix(text, "\n") {
		text = text[1:]
	}
	if strings.HasSuffix(text, "\n") {
		text = text[:len(text)-1]
	}
	return strings.Split(text, "\n")
}

func main() {
	rootDit := "example"
	challs := map[string](Challenge){}
	flags := map[string](Flag){}

	fmt.Println("=== Reading the challenges...")
	genres, err := os.ReadDir(rootDit)
	if err != nil {
		fmt.Printf("Error reading the directory %v: %v\n", ".", err)
		return
	}

	for _, genre := range genres {
		genrePath := fmt.Sprintf("%s/%s", rootDit, genre.Name())
		challDirs, err := os.ReadDir(genrePath)
		if err != nil {
			fmt.Printf("Error reading the directory %v: %v\n", ".", err)
			return
		}
		for _, chall := range challDirs {
			challPath := fmt.Sprintf("%s/%s", genrePath, chall.Name())
			files, err := os.ReadDir(challPath)
			if err != nil {
				fmt.Printf("Error reading the directory %v: %v\n", ".", err)
				return
			}
			for _, file := range files {
				if file.Name() != "challenge.yml" && file.Name() != "flag.txt" {
					continue
				}

				filePath := fmt.Sprintf("%s/%s", challPath, file.Name())
				fmt.Println("=== Reading the file: " + filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading the file %v: %v\n", file.Name(), err)
					return
				}
				if file.Name() == "challenge.yml" {
					parsed, err := ParseChall(filePath, content)
					if err != nil {
						continue
					}
					challs[filePath] = parsed
				}
				if file.Name() == "flag.txt" {
					parsed := ParseFlag(filePath, content)
					flags[filePath] = parsed
				}
			}
		}
	}
}
