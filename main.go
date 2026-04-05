package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadChalls(rootDir string, config Config) (map[string]Challenge, map[string]Flags, error) {
	challs := map[string]Challenge{}
	flags := map[string]Flags{}

	fmt.Println("== Loading the challenges...")
	genres := []string{}

	if config.Genre != nil {
		for _, genre := range config.Genre {
			genrePath := filepath.Join(rootDir, genre)
			if _, err := os.Stat(genrePath); os.IsNotExist(err) {
				fmt.Printf("%sWarning%s: the genre %v is not found.\n", colorYellow, colorReset, genre)
			} else {
				genres = append(genres, genre)
			}
		}
	} else {
		genreEnts, err := os.ReadDir(rootDir)
		if err != nil {
			fmt.Printf("%sError%s: loading the directory %v: %v\n", colorRed, colorReset, rootDir, err)
			return challs, flags, err
		}
		for _, genreEnt := range genreEnts {
			genres = append(genres, genreEnt.Name())
		}
	}

	for _, genre := range genres {
		genrePath := filepath.Join(rootDir, genre)
		challDirs, err := os.ReadDir(genrePath)
		if err != nil {
			fmt.Printf("%sError%s: loading the directory %v: %v\n", colorRed, colorReset, genrePath, err)
			return challs, flags, err
		}
		for _, chall := range challDirs {
			challPath := filepath.Join(genrePath, chall.Name())
			files, err := os.ReadDir(challPath)
			if err != nil {
				fmt.Printf("%sError%s: loading the directory %v: %v\n", colorRed, colorReset, challPath, err)
				return challs, flags, err
			}
			for _, file := range files {
				if file.Name() != "challenge.yml" && file.Name() != "challenge.yaml" && file.Name() != "flag.txt" {
					continue
				}

				filePath := filepath.Join(challPath, file.Name())
				fmt.Println("=== Loading the file: " + filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("%sError%s: loading the file %v: %v\n", colorRed, colorReset, file.Name(), err)
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
					parsed := ParseFlag(content)
					flags[challPath] = parsed
				}
			}
		}
	}
	return challs, flags, nil
}

func matchFlag(flag Flag, challFlag YamlFlag) bool {
	actually := flag.flag
	expected := challFlag.Content
	if challFlag.Data == "case_insensitive" {
		expected = strings.ToLower(expected)
		actually = strings.ToLower(actually)
	}
	switch challFlag.Type {
	case "static":
		if flag.fail {
			return actually != expected
		}
		return actually == expected
	case "regex":
		reg := regexp.MustCompile(expected)
		ms := reg.FindStringIndex(actually)
		return ms != nil && ms[0] == 0 && ms[1] == len(actually)
	}
	return false
}

func UnitTest(challs map[string]Challenge, flagMap map[string]Flags) bool {
	fmt.Println("== Unit testing...")
	isErr := false
	for challPath, flags := range flagMap {
		fmt.Printf("\n=== Unit testing the challenge: %v\n", challPath)
		chall := challs[challPath]

		if chall.Type != "" && chall.Type != "static" && chall.Type != "dynamic" {
			fmt.Printf("%sSKIP%s: Challenge skipped (type: %v) (%v)\n", colorYellow, colorReset, chall.Type, challPath)
			continue
		}

		for _, flag := range flags {
			ok := false
			for _, challFlag := range chall.Flags {
				if matchFlag(flag, challFlag) {
					ok = true
					break
				}
			}
			assertMsg := fmt.Sprintf("is assert_%sok%s", colorGreen, colorReset)
			if flag.fail {
				assertMsg = fmt.Sprintf("is assert_%sng%s", colorRed, colorReset)
			}
			if ok {
				fmt.Printf("%sPASS%s: %v %s (%v)\n", colorGreen, colorReset, flag.flag, assertMsg, challPath)
			} else {
				fmt.Printf("%sFAIL%s: %v %s (%v)\n", colorRed, colorReset, flag.flag, assertMsg, challPath)
				isErr = true
			}
		}
	}
	return isErr
}

func GetConfig(file string) (Config, error) {
	if file == "" {
		return Config{}, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("loading the config file %v: %w", file, err)
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshalling the config file %v: %w", file, err)
	}

	return config, nil
}

func main() {
	rootDir := os.Getenv("INPUT_TARGET_DIRECTORY")
	if rootDir == "" {
		rootDir = "."
	}

	configFile := os.Getenv("INPUT_CONFIG_FILE")

	config, err := GetConfig(configFile)
	if err != nil {
		fmt.Printf("%sError%s: %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	challs, flags, err := LoadChalls(rootDir, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("")
	isErr := UnitTest(challs, flags)
	if isErr {
		os.Exit(1)
	}
}
