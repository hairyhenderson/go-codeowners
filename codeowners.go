package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Codeowners - patterns/owners mappings for the given repo
type Codeowners struct {
	repoRoot string
	patterns []Codeowner
}

// Codeowner - owners for a given pattern
type Codeowner struct {
	pattern string
	ig      *gitignore.GitIgnore
	owners  []string
}

func (c Codeowner) String() string {
	return fmt.Sprintf("%s\t%v", c.pattern, strings.Join(c.owners, ", "))
}

var fs = afero.NewOsFs()

// findCodeownersFile - find a CODEOWNERS file somewhere within or below
// the working directory (wd), and open it.
func findCodeownersFile(wd string) (io.Reader, string, error) {
	dir := wd
	for {
		for _, p := range []string{".", "docs", ".github"} {
			pth := path.Join(dir, p)
			exists, err := afero.DirExists(fs, pth)
			if err != nil {
				return nil, "", err
			}
			if exists {
				f := path.Join(pth, "CODEOWNERS")
				_, err := fs.Stat(f)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return nil, "", err
				}
				r, err := fs.Open(f)
				return r, dir, err
			}
		}
		odir := dir
		dir = path.Dir(odir)
		// if we can't go up any further...
		if odir == dir {
			break
		}
		// if we're heading above the volume name (relevant on Windows)...
		if len(dir) < len(filepath.VolumeName(odir)) {
			break
		}
	}
	return nil, "", nil
}

// NewCodeowners -
func NewCodeowners(path string) (*Codeowners, error) {
	r, root, err := findCodeownersFile(path)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("No CODEOWNERS found in %s", path)
	}
	co := &Codeowners{
		repoRoot: root,
	}
	co.patterns, err = parseCodeowners(r)
	return co, nil
}

// parseCodeowners -
func parseCodeowners(r io.Reader) ([]Codeowner, error) {
	co := []Codeowner{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) > 0 && strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) > 1 {
			// for CODEOWNERS, * means all files, recursively
			// this differs from gitignore rules
			if fields[0] == "*" {
				fields[0] = "**"
			}
			c, err := NewCodeowner(fields[0], fields[1:])
			if err != nil {
				return nil, err
			}
			co = append(co, c)
		}
	}
	return co, nil
}

// NewCodeowner -
func NewCodeowner(pattern string, owners []string) (Codeowner, error) {
	ig, err := gitignore.CompileIgnoreLines(pattern)
	if err != nil {
		return Codeowner{}, err
	}
	c := Codeowner{
		pattern: pattern,
		ig:      ig,
		owners:  owners,
	}
	// fmt.Printf("NewCodeowner: %#v\n", c)
	return c, err
}

// Owners - return the list of code owners for the given path
// (within the repo root)
func (c *Codeowners) Owners(path string) []string {
	sort.Slice(c.patterns, sortPatterns(c.patterns))
	if strings.HasPrefix(path, c.repoRoot) {
		path = strings.Replace(path, c.repoRoot, "", 1)
	}
	for _, p := range c.patterns {
		if p.ig.MatchesPath(path) {
			return p.owners
		}
	}
	return nil
}

// returns a sort function to put the most-specific patterns first for
// consideration.
//
// e.g. foo/bar is more specific than foo/*
func sortPatterns(pats []Codeowner) func(i, j int) bool {
	return func(i, j int) bool {
		pi := pats[i].pattern
		pj := pats[j].pattern
		li := len(pi)
		lj := len(pj)
		if li > lj {
			return true
		}
		if li == lj {
			si := strings.Count(pi, "*")
			sj := strings.Count(pj, "*")
			return si < sj
		}
		return false
	}
}
