package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
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

	targetTime, err := time.Parse("2006-01-02", targetDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid date format: %v\n", err)
		os.Exit(1)
	}

	// Lipgloss styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	tagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	prevTagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	noneStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F")).Italic(true)

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

	// Print output
	for tagType := range tagPatterns {
		data := typeData[tagType]
		if data.Today != nil {
			fmt.Println(titleStyle.Render(fmt.Sprintf("%s tag on %s:", tagType, targetDate)),
				tagStyle.Render(data.Today.Name))
			if data.Previous != nil {
				fmt.Println("  Previous "+tagType+" tag:",
					prevTagStyle.Render(data.Previous.Name))
			} else {
				fmt.Println("  Previous "+tagType+" tag:", noneStyle.Render("(none)"))
			}
		} else {
			fmt.Println(noneStyle.Render("No " + tagType + " tag found on " + targetDate))
		}
	}
}
