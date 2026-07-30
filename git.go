package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var tagPatterns = map[string]*regexp.Regexp{
	"dev":  regexp.MustCompile(`^v[01]\.\d+\.\d+$`),
	"test": regexp.MustCompile(`^v[01]\.\d+\.\d+(-test\d+)?$`),
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
		if pattern.MatchString(tag) {
			matching = append(matching, tag)
		}
	}

	if len(matching) == 0 {
		return "", nil
	}

	best := matching[0]
	for _, tag := range matching[1:] {
		if compareTags(mode, tag, best) > 0 {
			best = tag
		}
	}
	return best, nil
}

func compareTags(mode, a, b string) int {
	switch mode {
	case "dev":
		ma, mia, pa, _ := parseSemVer(a)
		mb, mib, pb, _ := parseSemVer(b)
		if ma != mb {
			return ma - mb
		}
		if mia != mib {
			return mia - mib
		}
		return pa - pb
	case "test":
		ma, mia, pa, ta := parseSemVer(a)
		mb, mib, pb, tb := parseSemVer(b)
		if ma != mb {
			return ma - mb
		}
		if mia != mib {
			return mia - mib
		}
		if pa != pb {
			return pa - pb
		}
		return ta - tb
	case "uat", "prod", "img":
		da, sa := parseDateSeq(a)
		db, sb := parseDateSeq(b)
		if da != db {
			return da - db
		}
		return sa - sb
	default:
		return strings.Compare(a, b)
	}
}

var semVerRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-test(\d+))?$`)

func parseSemVer(tag string) (int, int, int, int) {
	matches := semVerRe.FindStringSubmatch(tag)
	if len(matches) < 4 {
		return 0, 0, 0, 0
	}
	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	testNum := 0
	if len(matches) >= 5 && matches[4] != "" {
		testNum, _ = strconv.Atoi(matches[4])
	}
	return major, minor, patch, testNum
}

var dateSeqRe = regexp.MustCompile(`(\d{8})\.(\d+)$`)

func parseDateSeq(tag string) (int, int) {
	matches := dateSeqRe.FindStringSubmatch(tag)
	if len(matches) != 3 {
		return 0, 0
	}
	date, _ := strconv.Atoi(matches[1])
	seq, _ := strconv.Atoi(matches[2])
	return date, seq
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
	// if err := os.MkdirAll(targetDir, 0755); err != nil {
	// fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", targetDir, err)
	// os.Exit(1)
	// }

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
