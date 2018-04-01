package codeowners

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"testing"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

var (
	sample = `# comment
*	@everyone

docs/**	@org/docteam @joe`
	sample2 = `* @hairyhenderson`
	sample3 = `baz/* @baz @qux`

	codeowners []Codeowner
)

func TestParseCodeowners(t *testing.T) {
	r := bytes.NewBufferString(sample)
	c, err := parseCodeowners(r)
	assert.NoError(t, err)
	expected := []Codeowner{
		co("**", []string{"@everyone"}),
		co("docs/**", []string{"@org/docteam", "@joe"}),
	}
	assert.Equal(t, expected, c)
}

func BenchmarkParseCodeowners(b *testing.B) {
	r := bytes.NewBufferString(sample)
	var c []Codeowner

	for n := 0; n < b.N; n++ {
		c, _ = parseCodeowners(r)
	}

	codeowners = c
}

func TestFindCodeownersFile(t *testing.T) {
	oldfs := fs
	defer func() {
		fs = oldfs
	}()
	fs = afero.NewMemMapFs()
	fs.Mkdir("/src/.github", 0755)
	fs.MkdirAll("/src/foo/bar/baz", 0755)
	fs.MkdirAll("/src/foo/qux/docs", 0755)
	fs.MkdirAll("/src/foo/qux/quux", 0755)
	f, _ := fs.Create("/src/.github/CODEOWNERS")
	f.WriteString(sample)

	f, _ = fs.Create("/src/foo/CODEOWNERS")
	f.WriteString(sample2)

	f, _ = fs.Create("/src/foo/qux/docs/CODEOWNERS")
	f.WriteString(sample3)

	r, root, err := findCodeownersFile("/src")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "/src", root)
	if r != nil {
		b, _ := ioutil.ReadAll(r)
		assert.Equal(t, sample, string(b))
	}

	r, root, err = findCodeownersFile("/src/foo/bar")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "/src/foo", root)
	if r != nil {
		b, _ := ioutil.ReadAll(r)
		assert.Equal(t, sample2, string(b))
	}

	r, root, err = findCodeownersFile("/src/foo/qux/quux")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "/src/foo/qux", root)
	if r != nil {
		b, _ := ioutil.ReadAll(r)
		assert.Equal(t, sample3, string(b))
	}

	r, root, err = findCodeownersFile("/")
	assert.NoError(t, err)
	assert.Nil(t, r)
}

func co(pattern string, owners []string) Codeowner {
	c, err := NewCodeowner(pattern, owners)
	if err != nil {
		panic(err)
	}
	return c
}

func TestOwners(t *testing.T) {
	foo := []string{"@foo"}
	bar := []string{"@bar"}
	baz := []string{"@baz"}
	data := []struct {
		patterns []Codeowner
		path     string
		expected []string
	}{
		{[]Codeowner{co("**", foo)}, "a/b", foo},
		{[]Codeowner{co("**", foo), co("a/b/*", bar)}, "a/b/c", bar},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/b/c", baz)}, "a/b/c", baz},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/*/c", baz)}, "a/b/c", bar},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/b/", baz)}, "a/b/bar", bar},
		{[]Codeowner{co("**", foo), co("/a/b/*", bar), co("a/b/", baz)}, "/someroot/a/b/bar", bar},
	}

	for _, d := range data {
		c := &Codeowners{patterns: d.patterns, repoRoot: "/someroot"}
		owners := c.Owners(d.path)
		assert.Equal(t, d.expected, owners)
	}
}

func cwd() string {
	_, filename, _, _ := runtime.Caller(0)
	cwd := path.Dir(filename)
	return cwd
}

func ExampleNewCodeowners() {
	c, _ := NewCodeowners(cwd())
	fmt.Println(c.patterns[0])
	// Output:
	// **	@hairyhenderson
}

func ExampleCodeowners_Owners() {
	c, _ := NewCodeowners(cwd())
	owners := c.Owners("README.md")
	for i, o := range owners {
		fmt.Printf("Owner #%d is %s\n", i, o)
	}
	// Output:
	// Owner #0 is @hairyhenderson
}
