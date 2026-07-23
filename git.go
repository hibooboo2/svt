package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hibooboo2/svt/exec"
)

var tagPatterns = map[string]*regexp.Regexp{
	"dev":  regexp.MustCompile(`^v\d+\.\d+\.\d+$`),
	"uat":  regexp.MustCompile(`^uat-\d{8}\.\d+$`),
	"prod": regexp.MustCompile(`^r\d{8}\.\d+$`),
	"img":  regexp.MustCompile(`^img-\d{8}\.\d+$`),
}

func gitOutput(args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "git %s failed: %v\n", strings.Join(args, " "), err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}

func gitRun(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "git %s failed: %v\n", strings.Join(args, " "), err)
		os.Exit(1)
	}
}

func gitFetch(tags bool) {
	args := []string{"fetch", "--all"}
	if tags {
		args = append(args, "--tags")
	}
	gitRun(args...)
}

func pull(branch string) {
	gitFetch(true)
	gitRun("checkout", branch)
	gitRun("pull")
}

func getCurrentBranch() string {
	return gitOutput("branch", "--show-current")
}

func findLastGitTag(mode string) (string, error) {
	if mode == "test" {
		mode = "dev"
	}
	pattern, ok := tagPatterns[mode]
	if !ok {
		return "", fmt.Errorf("unknown mode: %s", mode)
	}

	out, err := exec.Command("git", "tag", "--list").Output()
	if err != nil {
		return "", err
	}

	var matching []string
	for _, tag := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if mode == "dev" && strings.Contains(tag, "-test") {
			continue
		}
		if pattern.MatchString(tag) {
			matching = append(matching, tag)
		}
	}

	if len(matching) == 0 {
		return "", nil
	}

	sort.Strings(matching)
	return matching[len(matching)-1], nil
}

func tagIfNoTag(tag, branch string, pushArgs ...string) {
	currentBranch := getCurrentBranch()
	if currentBranch != branch {
		fmt.Fprintf(os.Stderr, "Error: Not on the '%s' branch.\n", branch)
		os.Exit(1)
	}

	if _, err := exec.Command("git", "rev-parse", tag).Output(); err == nil {
		fmt.Printf("Tag %s already exists globally.\n", tag)
		return
	}

	prevTag := decrementTagSuffix(tag)
	tagsOnCommit := gitOutput("tag", "--points-at", "HEAD")
	tagsList := strings.Split(tagsOnCommit, "\n")

	for _, t := range tagsList {
		t = strings.TrimSpace(t)
		if t == tag || (prevTag != "" && t == prevTag) {
			fmt.Printf("Commit already tagged with: %s or %s\n", tag, prevTag)
			return
		}
	}

	fmt.Printf("No matching tag found. Tagging with: %s\n", tag)
	gitRun("tag", tag)
	args := []string{"push", "origin", tag}
	args = append(args, pushArgs...)
	if err := exec.Command("git", args...).Run(); err != nil {
		gitRun("tag", "-d", tag)
		os.Exit(1)
	}
}

func decrementTagSuffix(tag string) string {
	re := regexp.MustCompile(`^(.*\D)(\d+)$`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) != 3 {
		return ""
	}
	prefix := matches[1]
	num := 0
	fmt.Sscanf(matches[2], "%d", &num)
	if num > 0 {
		return fmt.Sprintf("%s%d", prefix, num-1)
	}
	return ""
}

func symLinkBin() {
	binPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find binary: %v\n", err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find home dir: %v\n", err)
		os.Exit(1)
	}

	targetDir := filepath.Join(home, "go", "bin")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", targetDir, err)
		os.Exit(1)
	}

	names := []string{
		"ntag", "utag", "ptag", "itag", "ttag",
		"gntag", "gutag", "gptag", "gitag", "gttag", "gmtag",
		"gtags", "newMR", "symLinkBin",
	}

	for _, name := range names {
		linkPath := filepath.Join(targetDir, name)
		os.Remove(linkPath)
		if err := os.Symlink(binPath, linkPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to symlink %s: %v\n", linkPath, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Created symlinks in %s\n", targetDir)
	fmt.Printf("Ensure %s is in your PATH\n", targetDir)
}
