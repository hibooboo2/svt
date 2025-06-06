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

type TagSet struct {
	Previous *TagInfo
	Todays   []*TagInfo
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

	typeData := map[string]*TagSet{
		"r":   {},
		"v":   {},
		"uat": {},
	}

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
					typeData[tagType].Todays = append(typeData[tagType].Todays, &TagInfo{tag, tagTime})
				} else if tagTime.Before(targetTime) {
					prev := typeData[tagType].Previous
					if prev == nil || tagTime.After(prev.Date) {
						typeData[tagType].Previous = &TagInfo{tag, tagTime}
					}
				}
				break
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Git command failed: %v\n", err)
		os.Exit(1)
	}

	for _, tagType := range []string{"v", "uat", "r"} {
		data := typeData[tagType]
		if len(data.Todays) > 0 {
			latest := data.Todays[0]
			for _, t := range data.Todays[1:] {
				latest = compareTags(tagType, latest, t)
			}

			fmt.Println(titleStyle.Render(targetDate), "\t\t\t", tagStyle.Render(latest.Name))
			if data.Previous != nil {
				fmt.Println("Prev: ", int(targetTime.Truncate(time.Hour*24).Sub(data.Previous.Date.Truncate(time.Hour*24))/time.Hour/24), "Days Before\t\t", prevTagStyle.Render(data.Previous.Name))
			} else {
				fmt.Println("Prev: ", noneStyle.Render("(none)"))
			}
		} else {
			fmt.Println(noneStyle.Render("No " + tagType + " tag found on " + targetDate))
		}
	}
}

func compareTags(tagType string, a, b *TagInfo) *TagInfo {
	switch tagType {
	case "r", "uat":
		// Compare by the numeric suffix after the last dot
		getSuffix := func(tag string) int {
			parts := strings.Split(tag, ".")
			if len(parts) < 2 {
				return 0
			}
			var n int
			fmt.Sscanf(parts[len(parts)-1], "%d", &n)
			return n
		}
		if getSuffix(a.Name) >= getSuffix(b.Name) {
			return a
		}
		return b

	case "v":
		// Format: vX.Y.Z
		parseVer := func(tag string) (int, int, int) {
			var major, minor, patch int
			fmt.Sscanf(tag, "v%d.%d.%d", &major, &minor, &patch)
			return major, minor, patch
		}
		a1, a2, a3 := parseVer(a.Name)
		b1, b2, b3 := parseVer(b.Name)
		if a1 > b1 || (a1 == b1 && a2 > b2) || (a1 == b1 && a2 == b2 && a3 >= b3) {
			return a
		}
		return b
	}
	return a
}
