package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type TagInfo struct {
	Name string
	Date time.Time
}

var tagPatterns = map[string]*regexp.Regexp{
	"r":   regexp.MustCompile(`^r\d{8}\.\d+$`),
	"v":   regexp.MustCompile(`^v\d+\.\d+\.\d+$`),
	"uat": regexp.MustCompile(`^uat-\d{8}\.\d+$`),
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s YYYY-MM-DD\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s %s\n", os.Args[0], time.Now().Format("2006-01-02"))
		os.Args = append(os.Args, time.Now().Format("2006-01-02"))
	}
	targetDate := os.Args[1]

	// Parse input date
	targetTime, err := time.Parse("2006-01-02", targetDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid date format: %v\n", err)
		os.Exit(1)
	}

	// Run git command to get tags and creation dates
	cmd := exec.Command("git", "for-each-ref", "--sort=creatordate",
		"--format=%(refname:short) %(creatordate:iso)", "refs/tags")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get git output: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start git command: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(stdout)

	// Keep track of most recent and previous tags for each type
	typeData := map[string]struct {
		Previous *TagInfo
		Today    *TagInfo
	}{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		tag := parts[0]
		dateStr := parts[1]

		tagTime, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			continue
		}

		for tagType, pattern := range tagPatterns {
			if pattern.MatchString(tag) {
				if tagTime.Format("2006-01-02") == targetTime.Format("2006-01-02") {
					if typeData[tagType].Today == nil {
						data := typeData[tagType]
						data.Today = &TagInfo{tag, tagTime}
						typeData[tagType] = data
					}
				} else if tagTime.Before(targetTime) {
					data := typeData[tagType]
					data.Previous = &TagInfo{tag, tagTime}
					typeData[tagType] = data
				}
				break
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Git command failed: %v\n", err)
		os.Exit(1)
	}

	// Print results
	for tagType := range tagPatterns {
		data := typeData[tagType]
		if data.Today != nil {
			fmt.Printf("%s tag on %s: %s\n", tagType, targetDate, data.Today.Name)
			if data.Previous != nil {
				fmt.Printf("  Previous %s tag: %s\n", tagType, data.Previous.Name)
			} else {
				fmt.Printf("  Previous %s tag: (none)\n", tagType)
			}
		} else {
			fmt.Printf("No %s tag found on %s\n", tagType, targetDate)
		}
	}
}
