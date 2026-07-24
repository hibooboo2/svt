package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type binaryAction struct {
	mode       string
	gitWorkout bool
	branch     string
}

var binaryMap = map[string]binaryAction{
	"svl":        {},
	"ntag":       {mode: "dev"},
	"utag":       {mode: "uat"},
	"ptag":       {mode: "prod"},
	"itag":       {mode: "img"},
	"ttag":       {mode: "test"},
	"gntag":      {mode: "dev", gitWorkout: true, branch: "development"},
	"gutag":      {mode: "uat", gitWorkout: true, branch: "staging"},
	"gptag":      {mode: "prod", gitWorkout: true, branch: "production"},
	"gitag":      {mode: "img", gitWorkout: true},
	"gttag":      {mode: "test", gitWorkout: true},
	"gmtag":      {mode: "dev", gitWorkout: true, branch: "main"},
	"gtags":      {mode: "gtags"},
	"newMR":      {mode: "newMR"},
	"symLinkBin": {mode: "symLinkBin"},
}

func main() {
	binaryName := filepath.Base(os.Args[0])
	act, known := binaryMap[binaryName]
	if !known {
		symLinkBin()
		slog.Warn("Did not know bin, symlinking to all required ones", "bin", binaryName)
		return
	}

	defaultMode := "dev"
	if known && act.mode != "" {
		defaultMode = act.mode
	}

	lastTagFlag := flag.Bool("last-tag", false, "find the last tag from stdin")
	modeFlag := flag.String("mode", defaultMode, "version mode: dev, uat, prod, img, test")
	flag.Parse()

	mode := *modeFlag
	if mode == "" {
		mode = "dev"
	}

	if *lastTagFlag {
		handleLastTag(mode)
		return
	}

	if known && act.mode == "gtags" {
		handleGtags()
		return
	}
	if known && act.mode == "newMR" {
		handleNewMR()
		return
	}

	if known && act.gitWorkout {
		handleGitWorkout(mode, act.branch, flag.Args())
		return
	}

	handleTagGeneration(mode, flag.Args())
}

func handleTagGeneration(mode string, args []string) {
	if mode == "test" {
		if len(args) == 0 {
			gitFetch(true)
			last, _ := findLastGitTag("test")
			if last == "" {
				last = "v0.0.0"
			}
			args = []string{last}
		}
		fmt.Println(generateTestTag(args[0]))
		return
	}

	if len(args) == 0 {
		gitFetch(true)
		last, err := findLastGitTag(mode)
		if err != nil || last == "" {
			last = defaultTag(mode)
		}
		args = []string{last}
		if mode == "dev" {
			args = append(args, "v0.0.1")
		}
	}

	tag := computeNextTag(mode, args)
	fmt.Println(tag)
}

func handleGitWorkout(mode string, branch string, pushArgs []string) {
	if branch == "" {
		branch = getCurrentBranch()
	}

	pull(branch)
	gitFetch(true)

	last, err := findLastGitTag(mode)
	if err != nil || last == "" {
		last = defaultTag(mode)
	}

	tagArgs := []string{last}
	if mode == "dev" {
		tagArgs = append(tagArgs, "v0.0.1")
	}

	tag := computeNextTag(mode, tagArgs)
	fmt.Println(tag)
	tagIfNoTag(tag, branch, pushArgs...)
}

func computeNextTag(mode string, args []string) string {
	if mode == "test" {
		if len(args) == 0 {
			return generateTestTag("v0.0.0")
		}
		return generateTestTag(args[0])
	}

	var v Version
	var v2 Version
	switch mode {
	case "dev":
		if len(args) < 2 {
			v = SemVer(args[0])
			v2 = SemVer("v0.0.1")
		} else {
			v = SemVer(args[0])
			v2 = SemVer(args[1])
		}
	case "uat":
		v = UAT(args[0])
	case "prod", "img":
		v = PROD(args[0])
	default:
		fmt.Fprintf(os.Stderr, "invalid mode: %s\n", mode)
		os.Exit(1)
	}
	v = v.Version(v2)
	switch mode {
	case "img":
		s := fmt.Sprintf("%s", v)
		return "img-" + s[1:]
	default:
		return fmt.Sprintf("%s", v)
	}
}

func generateTestTag(baseTag string) string {
	if strings.Contains(baseTag, "-test") {
		idx := strings.Index(baseTag, "-test")
		mainTag := baseTag[:idx]
		rest := baseTag[idx+5:]
		num := 0
		if rest != "" {
			fmt.Sscanf(rest, "%d", &num)
		}
		return fmt.Sprintf("%s-test%d", mainTag, num+1)
	}
	return baseTag + "-test1"
}

func defaultTag(mode string) string {
	switch mode {
	case "dev", "test":
		return "v0.0.0"
	case "uat":
		return fmt.Sprintf("uat-%s.0", dateStringUATFormat(time.Now()))
	case "prod", "img":
		return fmt.Sprintf("r%s.0", dateStringUATFormat(time.Now()))
	}
	return ""
}

func handleLastTag(mode string) {
	scanner := bufio.NewScanner(os.Stdin)
	var tags []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		tags = append(tags, line)
	}

	if len(tags) == 0 {
		return
	}

	best := tags[0]
	for _, tag := range tags[1:] {
		if compareTags(mode, tag, best) > 0 {
			best = tag
		}
	}
	fmt.Println(best)
}

func handleGtags() {
	out, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		return
	}
	parts := strings.SplitN(string(out), "-", 2)
	if len(parts) == 2 {
		fmt.Println(strings.TrimSpace(parts[0]) + "-" + strings.SplitN(parts[1], "-", 2)[0])
	} else {
		fmt.Print(strings.TrimSpace(string(out)))
	}
}

func handleNewMR() {
	cwd, _ := os.Getwd()
	repo := filepath.Base(cwd)
	branch := getCurrentBranch()
	fmt.Printf("https://git.nops.ftr.com/raven/services/%s/-/merge_requests/new?merge_request%%5Bsource_branch%%5D=%s\n", repo, branch)
}

type Version interface {
	Version(...Version) Version
}
