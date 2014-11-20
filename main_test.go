package main

import (
	"io/ioutil"
	logPkg "log"
	"os"
	"testing"

	assert "github.com/pilu/miniassert"
)

func init() {
	log.verbose = false
}

func TestNewPlate(t *testing.T) {
	p := newPlate("foo", "bar")
	assert.Equal(t, "foo", p.srcPath)
	assert.Equal(t, "bar", p.outPath)
}

func TestPlate_BuildTemplatePath(t *testing.T) {
	p := newPlate("foo", "bar")
	assert.Equal(t, "foo/go.tpl", p.buildTemplatePath("go"))
}

func TestPlate_BuildOutPath(t *testing.T) {
	p := newPlate("foo", "bar")
	assert.Equal(t, "bar/tmp/file.go", p.buildOutPath("tmp/file.go"))
}

func TestPlate_OpenTemplate(t *testing.T) {
	p := newPlate("__test-fixtures__", "tmp")
	mainTpl, err := p.openTemplate("go")
	assert.NotNil(t, mainTpl)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(mainTpl.Templates()))

	names := []string{"main.go", "main_test.go", "config/app.config"}
	for _, name := range names {
		tpl := mainTpl.Lookup(name)
		assert.NotNil(t, tpl)
	}
}

func TestPlate_OpenTemplate_template_not_found(t *testing.T) {
	p := newPlate("__test-fixtures__", "tmp")
	_, err := p.openTemplate("random-name")
	assert.NotNil(t, err)
}

func TestPlate_Execute(t *testing.T) {
	defer os.RemoveAll("__test-data__")
	p := newPlate("__test-fixtures__", "__test-data__")
	p.execute("go", []string{"program-name", "outpath", "World"}...)

	paths := map[string]string{
		"__test-data__/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}`,
		"__test-data__/main_test.go": `package main

import "testing"

func TestFoo(t *testing.T) {
}`,
		"__test-data__/config/app.config": `config file`,
	}

	for path, expectedContent := range paths {
		f, err := os.Open(path)
		assert.Nil(t, err)

		content, err := ioutil.ReadAll(f)
		if err != nil {
			logPkg.Fatal(err)
		}
		assert.Equal(t, expectedContent, string(content))
	}
}
