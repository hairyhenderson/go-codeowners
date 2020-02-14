package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

// Codeowners - patterns/owners mappings for the given repo
type Codeowners struct {
	repoRoot string
	patterns []Codeowner
}

// Codeowner - owners for a given pattern
type Codeowner struct {
	pattern string
	re      *regexp.Regexp
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
		for _, p := range []string{".", "docs", ".github", ".gitlab"} {
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
	co.patterns = parseCodeowners(r)
	return co, nil
}

// parseCodeowners parses a list of Codeowners from a Reader
func parseCodeowners(r io.Reader) []Codeowner {
	co := []Codeowner{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) > 0 && strings.HasPrefix(fields[0], "#") {
			continue
		}
		if len(fields) > 1 {
			c, _ := NewCodeowner(fields[0], fields[1:])
			co = append(co, c)
		}
	}
	return co
}

// NewCodeowner -
func NewCodeowner(pattern string, owners []string) (Codeowner, error) {
	re := getPattern(pattern)
	c := Codeowner{
		pattern: pattern,
		re:      re,
		owners:  owners,
	}
	return c, nil
}

// Owners - return the list of code owners for the given path
// (within the repo root)
func (c *Codeowners) Owners(path string) []string {
	if strings.HasPrefix(path, c.repoRoot) {
		path = strings.Replace(path, c.repoRoot, "", 1)
	}

	// Order is important; the last matching pattern takes the most precedence.
	for i := len(c.patterns) - 1; i >= 0; i-- {
		p := c.patterns[i]

		if p.re.MatchString(path) {
			return p.owners
		}
	}

	return nil
}

// based on github.com/sabhiram/go-gitignore
// but modified so that 'dir/*' only matches files in 'dir/'
func getPattern(line string) *regexp.Regexp {
	// when # or ! is escaped with a \
	if regexp.MustCompile(`^(\\#|\\!)`).MatchString(line) {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern
}
