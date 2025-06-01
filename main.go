package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadChalls(rootDir string, config Config) (map[string](Challenge), map[string](Flags), error) {
	challs := map[string](Challenge){}
	flags := map[string](Flags){}

	fmt.Println("== Loading the challenges...")
	genres := []string{}

	if config.Genre != nil {
		for _, genre := range config.Genre {
			genrePath := fmt.Sprintf("%s/%s", rootDir, genre)
			if _, err := os.Stat(genrePath); os.IsNotExist(err) {
				fmt.Printf("\x1b[33mWarning\x1b[0m: the genre %v is not found.\n", genre)
			} else {
				genres = append(genres, genre)
			}
		}
	} else {
		genreEnts, err := os.ReadDir(rootDir)
		if err != nil {
			fmt.Printf("\x1b[31mError\x1b[0m: loading the directory %v: %v\n", rootDir, err)
			return challs, flags, err
		}
		for _, genreEnt := range genreEnts {
			genres = append(genres, genreEnt.Name())
		}
	}

	for _, genre := range genres {
		genrePath := fmt.Sprintf("%s/%s", rootDir, genre)
		challDirs, err := os.ReadDir(genrePath)
		if err != nil {
			fmt.Printf("\x1b[31mError\x1b[0m: loading the directory %v: %v\n", ".", err)
			return challs, flags, err
		}
		for _, chall := range challDirs {
			challPath := fmt.Sprintf("%s/%s", genrePath, chall.Name())
			files, err := os.ReadDir(challPath)
			if err != nil {
				fmt.Printf("\x1b[31mError\x1b[0m: loading the directory %v: %v\n", ".", err)
				return challs, flags, err
			}
			for _, file := range files {
				if file.Name() != "challenge.yml" && file.Name() != "challenge.yaml" && file.Name() != "flag.txt" {
					continue
				}

				filePath := fmt.Sprintf("%s/%s", challPath, file.Name())
				fmt.Println("=== Leading the file: " + filePath)

				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("\x1b[31mError\x1b[0m: loading the file %v: %v\n", file.Name(), err)
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

func UnitTest(challs map[string](Challenge), flagMap map[string](Flags)) bool {
	fmt.Println("== Unit testing...")
	isErr := false
	for challPath, flags := range flagMap {
		fmt.Printf("\n=== Unit testing the challenge: %v\n", challPath)
		chall := challs[challPath]

		// Skip testing if challenge type is neither static nor dynamic
		if chall.Type != "" && chall.Type != "static" && chall.Type != "dynamic" {
			fmt.Printf("\x1b[33mSKIP\x1b[0m: Challenge skipped (type: %v) (%v)\n", chall.Type, challPath)
			continue
		}

		for _, flag := range flags {
			ok := false
			for _, challFlags := range chall.Flags {
				actually := flag.flag
				fail := flag.fail
				expected := challFlags.Content
				caseInSensitive := challFlags.Data == "case_insensitive"
				if caseInSensitive {
					expected = strings.ToLower(expected)
					actually = strings.ToLower(actually)
				}
				if challFlags.Type == "static" {
					if fail && actually != expected {
						ok = true
						break
					} else if !fail && actually == expected {
						ok = true
						break
					}
				}
				if challFlags.Type == "regex" {
					reg := regexp.MustCompile(expected)
					ms := reg.FindStringIndex(actually)
					if ms != nil && ms[0] == 0 && ms[1] == len(actually) {
						ok = true
						break
					}
				}
			}
			assertMsg := "is assert_\x1b[32mok\x1b[0m"
			if flag.fail {
				assertMsg = "is assert_\x1b[31mng\x1b[0m"
			}
			if ok {
				fmt.Printf("\x1b[32mPASS\x1b[0m: %v %s (%v)\n", flag.flag, assertMsg, challPath)
			} else {
				fmt.Printf("\x1b[31mFAIL\x1b[0m: %v %s (%v)\n", flag.flag, assertMsg, challPath)
				isErr = true
			}
		}
	}
	return isErr
}

func GetConfig(file string) Config {
	if file == "" {
		return Config{}
	}

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("\x1b[31mError\x1b[0m: loading the config file %v: %v\n", file, err)
		return Config{}
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		fmt.Printf("\x1b[31mError\x1b[0m: unmarshalling the file %v: %v\n", file, err)
		return Config{}
	}

	return config
}

func main() {
	rootDir := os.Getenv("INPUT_TARGET_DIRECTORY")
	if rootDir == "" {
		rootDir = "."
	}

	configFile := os.Getenv("INPUT_CONFIG_FILE")

	config := GetConfig(configFile)

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
