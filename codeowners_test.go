package codeowners

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	sample = `# comment
*	@everyone

   foobar/  someone@else.com # inline comment

docs/**	@org/docteam @joe`
	sample2 = `* @hairyhenderson`
	sample3 = `baz/* @baz @qux`
	sample4 = `[test]
*   @everyone
[test2][2]
*/foo @everyoneelse`

	// based on https://help.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners#codeowners-syntax
	// with a few unimportant modifications
	fullSample = `# This is a comment.
# Each line is a file pattern followed by one or more owners.

# These owners will be the default owners for everything in
# the repo. Unless a later match takes precedence,
# @global-owner1 and @global-owner2 will be requested for
# review when someone opens a pull request.
*       @global-owner1 @global-owner2

# Order is important; the last matching pattern takes the most
# precedence. When someone opens a pull request that only
# modifies JS files, only @js-owner and not the global
# owner(s) will be requested for a review.
*.js	@js-owner

# You can also use email addresses if you prefer. They'll be
# used to look up users just like we do for commit author
# emails.
*.go docs@example.com

# In this example, @doctocat owns any files in the build/logs
# directory at the root of the repository and any of its
# subdirectories.
/build/logs/ @doctocat

# In this example, @fooowner owns any files in the /cells/foo
# directory at the root of the repository and any of its
# subdirectories and files.
/cells/foo @fooowner

# The 'docs/*' pattern will match files like
# 'docs/getting-started.md' but not further nested files like
# 'docs/build-app/troubleshooting.md'.
docs/*  docs@example.com

# In this example, @octocat owns any file in an apps directory
# anywhere in your repository.
apps/ @octocat

# In this example, @doctocat owns any file in the '/docs'
# directory in the root of your repository.
/docs/ @doctocat

  foobar/ @fooowner

\#foo/ @hashowner

docs/*.md @mdowner

# this example tests an escaped space in the path
space/test\ space/ @spaceowner

# In this example, @infra owns any file and directory in the
# '/terraform' directory in the root of your repository.
/terraform @infra
`

	gitlabSections = `# This is a GitLab section with default owners.
[Team 1][1] @default1 @default2
*.js	@js-owner
*.txt

# This is another section with new defaults.
[Team 2] @default3 @default4
*.java @java-owner
*

# This is an optional sections without defaults.
^[Team 3]
*.go @team-1
`

	codeowners []Codeowner
)

func TestParseGitLabSectionsWithDefaults(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(gitlabSections)
	c, _ := parseCodeowners(r)
	expected := []Codeowner{
		co("*.js", []string{"@js-owner"}),
		co("*.txt", []string{"@default1", "@default2"}),
		co("*.java", []string{"@java-owner"}),
		co("*", []string{"@default3", "@default4"}),
		co("*.go", []string{"@team-1"}),
	}
	assert.Equal(t, expected, c)
}

func TestParseCodeowners(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(sample)
	c, _ := parseCodeowners(r)
	expected := []Codeowner{
		co("*", []string{"@everyone"}),
		co("foobar/", []string{"someone@else.com"}),
		co("docs/**", []string{"@org/docteam", "@joe"}),
	}
	assert.Equal(t, expected, c)
}

func TestParseCodeownersSections(t *testing.T) {
	t.Parallel()
	r := bytes.NewBufferString(sample4)
	c, _ := parseCodeowners(r)
	expected := []Codeowner{
		co("*", []string{"@everyone"}),
		co("*/foo", []string{"@everyoneelse"}),
	}
	assert.Equal(t, expected, c)
}

func BenchmarkParseCodeowners(b *testing.B) {
	var c []Codeowner

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		r := bytes.NewBufferString(sample)
		b.StartTimer()
		c, _ = parseCodeowners(r)
	}

	codeowners = c
}

func BenchmarkOwners(b *testing.B) {
	c, _ := FromReader(strings.NewReader(fullSample), "")
	data := []string{
		"#foo/bar.go",
		"blah/docs/README.md",
		"foo/bar/docs/foo/foo.js",
		"/space/test space/doc1.txt",
		"/terraform/kubernetes",
	}

	for _, d := range data {
		b.Run(d, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = c.Owners(d)
			}
		})
	}
}

func TestFindCodeownersFile(t *testing.T) {
	fsys := fstest.MapFS{
		"src/.github/CODEOWNERS":      &fstest.MapFile{Data: []byte(sample)},
		"src/foo/CODEOWNERS":          &fstest.MapFile{Data: []byte(sample2)},
		"src/foo/qux/docs/CODEOWNERS": &fstest.MapFile{Data: []byte(sample3)},
		"src/bar/CODEOWNERS":          &fstest.MapFile{Data: []byte(sample2)},
		"src/bar/.github/CODEOWNERS":  &fstest.MapFile{Data: []byte(sample)},
		"src/bar/.gitlab/CODEOWNERS":  &fstest.MapFile{Data: []byte(sample3)},
	}

	r, root, err := findCodeownersFile(fsys, "src")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src", root)

	b, _ := io.ReadAll(r)
	assert.Equal(t, sample, string(b))

	r, root, err = findCodeownersFile(fsys, "src/foo/bar")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src/foo", root)

	b, _ = io.ReadAll(r)
	assert.Equal(t, sample2, string(b))

	r, root, err = findCodeownersFile(fsys, "src/foo/qux/quux")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src/foo/qux", root)

	b, _ = io.ReadAll(r)
	assert.Equal(t, sample3, string(b))

	r, _, err = findCodeownersFile(fsys, ".")
	require.NoError(t, err)
	assert.Nil(t, r)

	r, root, err = findCodeownersFile(fsys, "src/bar")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "src/bar", root)

	b, _ = io.ReadAll(r)
	assert.Equal(t, sample, string(b))
}

func co(pattern string, owners []string) Codeowner {
	c, err := NewCodeowner(pattern, owners)
	if err != nil {
		panic(err)
	}
	return c
}

func TestFullParseCodeowners(t *testing.T) {
	t.Parallel()

	c, _ := parseCodeowners(strings.NewReader(fullSample))
	codeowners := &Codeowners{
		repoRoot: "/build",
		Patterns: c,
	}

	// these tests were ported from https://github.com/softprops/codeowners
	data := []struct {
		path   string
		owners []string
	}{
		{"#foo/bar.go", []string{"@hashowner"}},
		{"foobar/baz.go", []string{"@fooowner"}},
		{"/docs/README.md", []string{"@mdowner"}},
		// XXX: uncertain about this one
		{"blah/docs/README.md", []string{"docs@example.com"}},
		{"foo.txt", []string{"@global-owner1", "@global-owner2"}},
		{"foo/bar.txt", []string{"@global-owner1", "@global-owner2"}},
		{"foo.js", []string{"@js-owner"}},
		{"foo/bar.js", []string{"@js-owner"}},
		{"foo.go", []string{"docs@example.com"}},
		{"foo/bar.go", []string{"docs@example.com"}},
		// relative to root
		{"build/logs/foo.go", []string{"@doctocat"}},
		{"build/logs/foo/bar.go", []string{"@doctocat"}},
		// not relative to root
		{"foo/build/logs/foo.go", []string{"docs@example.com"}},
		// docs anywhere
		{"foo/docs/foo.js", []string{"docs@example.com"}},
		{"foo/bar/docs/foo.js", []string{"docs@example.com"}},
		// but not nested
		{"foo/bar/docs/foo/foo.js", []string{"@js-owner"}},
		{"foo/apps/foo.js", []string{"@octocat"}},
		{"docs/foo.js", []string{"@doctocat"}},
		{"/docs/foo.js", []string{"@doctocat"}},
		{"/space/test space/doc1.txt", []string{"@spaceowner"}},
		{"/terraform/kubernetes", []string{"@infra"}},

		{"/cells/foo", []string{"@fooowner"}},
		{"/cells/foo/", []string{"@fooowner"}},
		{"/cells/foo/bar", []string{"@fooowner"}},
		{"/cells/foo/bar/quux", []string{"@fooowner"}},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%q==%#v", d.path, d.owners), func(t *testing.T) {
			assert.EqualValues(t, d.owners, codeowners.Owners(d.path))
		})
	}
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
		{[]Codeowner{co("a/*", foo)}, "c/b", nil},
		{[]Codeowner{co("**", foo)}, "a/b", foo},
		{[]Codeowner{co("**", foo), co("a/b/*", bar)}, "a/b/c", bar},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/b/c", baz)}, "a/b/c", baz},
		{[]Codeowner{co("**", foo), co("a/*/c", bar), co("a/b/*", baz)}, "a/b/c", baz},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/b/", baz)}, "a/b/bar", baz},
		{[]Codeowner{co("**", foo), co("a/b/*", bar), co("a/b/", baz)}, "/someroot/a/b/bar", baz},
		{[]Codeowner{
			co("*", foo),
			co("/a/*", bar),
			co("/b/**", baz)}, "/a/aa/file", foo},
		{[]Codeowner{
			co("*", foo),
			co("/a/**", bar)}, "/a/bb/file", bar},
		{[]Codeowner{
			co("*", []string{"@foo", "@bar"}),
			co("/bar/", bar)}, "/bar/quux", bar},
		{[]Codeowner{
			co("*", []string{"@foo", "@bar"}),
			co("/bar", bar)}, "/bar", bar},
		{[]Codeowner{
			co("*", []string{"@foo", "@bar"}),
			co("/bar", bar)}, "/bar/quux", bar},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%s==%s", d.path, d.expected), func(t *testing.T) {
			c := &Codeowners{Patterns: d.patterns, repoRoot: "/someroot"}
			owners := c.Owners(d.path)
			assert.Equal(t, d.expected, owners)
		})
	}
}

func TestCombineEscapedSpaces(t *testing.T) {
	data := []struct {
		fields   []string
		expected []string
	}{
		{[]string{"docs/", "@owner"}, []string{"docs/", "@owner"}},
		{[]string{"docs/bob/**", "@owner"}, []string{"docs/bob/**", "@owner"}},
		{[]string{"docs/bob\\", "test/", "@owner"}, []string{"docs/bob test/", "@owner"}},
		{[]string{"docs/bob\\", "test/sub/final\\", "space/", "@owner"}, []string{"docs/bob test/sub/final space/", "@owner"}},
		{[]string{"docs/bob\\", "test/another\\", "test/**", "@owner"}, []string{"docs/bob test/another test/**", "@owner"}},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%s==%s", d.fields, d.expected), func(t *testing.T) {
			assert.Equal(t, d.expected, combineEscapedSpaces(d.fields))
		})
	}
}

func cwd() string {
	_, filename, _, _ := runtime.Caller(0)
	cwd := path.Dir(filename)
	return cwd
}

func ExampleFromFile() {
	c, _ := FromFile(cwd())
	fmt.Println(c.Patterns[0])
	// Output:
	// *	@hairyhenderson
}

func ExampleFromFileWithFS() {
	// open filesystem rooted at current working directory
	fsys := os.DirFS(cwd())

	c, _ := FromFileWithFS(fsys, ".")
	fmt.Println(c.Patterns[0])
	// Output:
	// *	@hairyhenderson
}

func ExampleFromReader() {
	reader := strings.NewReader(sample2)
	c, _ := FromReader(reader, "")
	fmt.Println(c.Patterns[0])
	// Output:
	// *	@hairyhenderson
}

func ExampleCodeowners_Owners() {
	c, _ := FromFile(cwd())
	owners := c.Owners("README.md")
	for i, o := range owners {
		fmt.Printf("Owner #%d is %s\n", i, o)
	}
	// Output:
	// Owner #0 is @hairyhenderson
}
