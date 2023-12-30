package main

import "fmt"

func func1(a int, b int) int {
	return a + b;
}

func main() {
	rootDit := "example"
	// var challenges []Challenge
	flags := map[string]([]string){}

	genres, err := os.ReadDir(rootDit)
	if err != nil {
		fmt.Printf("Error reading the directory %v: %v\n", ".", err)
		return
	}

	for _, genre := range genres {
		genrePath := fmt.Sprintf("%s/%s", rootDit, genre.Name())
		challs, err := os.ReadDir(genrePath)
		if err != nil {
			fmt.Printf("Error reading the directory %v: %v\n", ".", err)
			return
		}
		for _, chall := range challs {
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
				content, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Error reading the file %v: %v\n", file.Name(), err)
					return
				}
				if file.Name() == "challenge.yml" {
					fmt.Println(file.Name())
				}
				if file.Name() == "flag.txt" {
					parsed := ParseFlag(filePath, content)
					flags[parsed.Id] = parsed.Flags
				}
			}
		}
	}
}
