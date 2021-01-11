package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type DelimiterLineRange struct {
	Start int `json:"start,omitempty"`
	Stop  int `json:"stop,omitempty"`
}

// ReadFile reades files and writes
func ReadFile(filename string) []string {

	var (
		slice              = []string{}
		startConcatenation = false
		temp               string
		outputPrefix       = "-"
	)

	if Configuration.Commit.Output != nil {
		if Configuration.Commit.Output.Prefix != "" {
			outputPrefix = Configuration.Commit.Output.Prefix
		}
	}

	// opens file with read and write permissions
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(filename)
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Will close the file after main function is finished
	defer file.Close()

	// Grab contents of file
	read, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// File content to string
	newContents := string(read)
	lines := strings.Split(newContents, "\n")

	modifyIndexes := map[int]*DelimiterLineRange{}

	// Setting a scanner buffer layer
	scanner := bufio.NewScanner(file)

	// loop through the file.
	lineCount := -1

	currentModifiedIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		for i := 0; i < len(line); i++ {
			// Identifies a start of a delimiter
			if StartDelimiter(line, i) && startConcatenation == false {
				startConcatenation = true
				modifyIndexes[currentModifiedIndex] = &DelimiterLineRange{
					Start: lineCount,
				}

				temp = outputPrefix + " "
				i++
			} else {
				if startConcatenation == true {
					// Identifies if this is the second start delimiter before a end is found.
					if StartDelimiter(line, i) && startConcatenation == true {
						temp = ""
						startConcatenation = true
						modifyIndexes[currentModifiedIndex] = &DelimiterLineRange{
							Start: lineCount,
						}
						i++
					} else if EndDelimiter(line, i) {
						startConcatenation = false

						if obj, ok := modifyIndexes[currentModifiedIndex]; ok {
							obj.Stop = lineCount
							currentModifiedIndex++
						}

						slice = append(slice, temp)
						temp = outputPrefix + " "
						i++
					} else {
						temp += string(line[i])
					}
				}
			}
		}
		// Add space between words when the end delimiter doesnt end on the following line.
		temp += " "
	}

	deleteLines := []int{}
	for _, delimitersLocation := range modifyIndexes {
		for i, val := range lines {
			if i >= delimitersLocation.Start && i <= delimitersLocation.Stop {
				// TODO: If commit.remove is enabled remove the text.
				if Configuration.Commit.RemoveText {
					startIndex := strings.Index(val, "(:")
					endIndex := strings.Index(val, ":)")

					// Both delimiters are present
					if startIndex != -1 && endIndex != -1 {
						lines[i] = val[:startIndex] + val[endIndex+2:]
					}

					// no delimiter is present
					if startIndex == -1 && endIndex == -1 {
						// lines[i] = ""                // Erlinesse llinesst element (write zero vlineslue).
						// lines = lines[:len(lines)-1] // Trunclineste slice.
						deleteLines = append(deleteLines, i)
						r := strings.NewReplacer(val, "")
						lines[i] = r.Replace(val)
						continue
					}

					if startIndex != -1 && endIndex == -1 {
						lines[i] = val[:startIndex]
					}

					if endIndex != -1 && startIndex == -1 {
						lines[i] = val[endIndex+2:]
					}

					if lines[i] == "" {
						deleteLines = append(deleteLines, i)
					}

				} else {
					if strings.Contains(val, "(:") || strings.Contains(val, ":)") {
						r := strings.NewReplacer("(:", "", ":)", "")
						lines[i] = r.Replace(val)
					}
				}
			}
		}
	}

	// Remove the empty lines from removing the commit text if enabled.
	if len(deleteLines) > 0 {
		k := 0
		for index := 0; index < len(lines); index++ {
			not := false
			for _, line := range deleteLines {
				if index == line {
					not = true
				}
			}

			if !not {
				lines[k] = lines[index]
				k++
			}
		}

		lines = lines[:k]
	}

	// Join lines slice on newline
	newContents = strings.Join(lines, "\n")

	err = ioutil.WriteFile(filename, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}

	// save changes
	file.Sync()
	return slice
}

func StartDelimiter(line string, index int) bool {
	if string(line[index]) == "(" && index+1 < len(line) && string(line[index+1]) == ":" {
		return true
	}
	return false
}

func EndDelimiter(line string, index int) bool {
	if string(line[index]) == ":" && index+1 < len(line) && string(line[index+1]) == ")" {
		return true
	}
	return false
}
