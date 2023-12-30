package main

import (
	"fmt"
	"os"
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
				return chall, fmt.Errorf("\x1b[31mError\x1b[0m: flag content is not specified")
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
			return chall, fmt.Errorf("\x1b[31mError\x1b[0m: unknown flag type: %v", flag)
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

func LoadChalls(rootDir string) (map[string](Challenge), map[string](Flag), error) {
	challs := map[string](Challenge){}
	flags := map[string](Flag){}

	fmt.Println("== Reading the challenges...")
	genres, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Printf("\x1b[31mError\x1b[0m: reading the directory %v: %v\n", ".", err)
		return challs, flags, err
	}

	for _, genre := range genres {
		genrePath := fmt.Sprintf("%s/%s", rootDir, genre.Name())
		challDirs, err := os.ReadDir(genrePath)
		if err != nil {
			fmt.Printf("\x1b[31mError\x1b[0m: reading the directory %v: %v\n", ".", err)
			return challs, flags, err
		}
		for _, chall := range challDirs {
			challPath := fmt.Sprintf("%s/%s", genrePath, chall.Name())
			files, err := os.ReadDir(challPath)
			if err != nil {
				fmt.Printf("\x1b[31mError\x1b[0m: reading the directory %v: %v\n", ".", err)
				return challs, flags, err
			}
			for _, file := range files {
				if file.Name() != "challenge.yml" && file.Name() != "challenge.yaml" && file.Name() != "flag.txt" {
					continue
				}

				filePath := fmt.Sprintf("%s/%s", challPath, file.Name())
				fmt.Println("=== Reading the file: " + filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("\x1b[31mError\x1b[0m: reading the file %v: %v\n", file.Name(), err)
					return challs, flags, err
				}
				if file.Name() == "challenge.yml" || file.Name() == "challenge.yaml" {
					parsed, err := ParseChall(filePath, content)
					if err != nil {
						continue
					}
					challs[challPath] = parsed
				}
				if file.Name() == "flag.txt" {
					parsed := ParseFlag(filePath, content)
					flags[challPath] = parsed
				}
			}
		}
	}
	return challs, flags, nil
}

func UnitTest(challs map[string](Challenge), flagMap map[string](Flag)) bool {
	fmt.Println("== Unit testing...")
	isErr := false
	for challPath, flags := range flagMap {
		fmt.Printf("\n=== Unit testing the challenge: %v\n", challPath)
		chall := challs[challPath]
		for _, flag := range flags {
			ok := false
			for _, challFlags := range chall.Flags {
				actually := flag
				expected := challFlags.Content
				caseInSensitive := challFlags.Data == "case_insensitive"
				if caseInSensitive {
					expected = strings.ToLower(expected)
					actually = strings.ToLower(actually)
				}
				if challFlags.Type == "static" && actually == expected {
					ok = true
					break
				}
				if challFlags.Type == "regex" {
					reg := regexp.MustCompile(expected)
					if reg.MatchString(actually) {
						ok = true
						break
					}
				}
			}
			if ok {
				fmt.Printf("\x1b[32mPASS\x1b[0m: %v (%v)\n", flag, challPath)
			} else {
				fmt.Printf("\x1b[31mFAIL\x1b[0m: %v (%v)\n", flag, challPath)
				isErr = true
			}
		}
	}
	return isErr
}

func main() {
	rootDir := os.Getenv("INPUT_TARGET_DIRECTORY")
	if rootDir == "" {
		rootDir = "."
	}

	challs, flags, err := LoadChalls(rootDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("")
	isErr := UnitTest(challs, flags)
	if isErr {
		os.Exit(1)
	}
}
