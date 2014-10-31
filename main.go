package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templatesExtension = ".tpl"
	templatesFolder    = ".plate"
)

type plate struct {
	srcPath string
	outPath string
}

func newPlate(srcPath, outPath string) *plate {
	return &plate{
		srcPath: srcPath,
		outPath: outPath,
	}
}

func (p *plate) buildTemplatePath(name string) string {
	filename := fmt.Sprintf("%s%s", name, templatesExtension)
	return path.Join(p.srcPath, filename)
}

func (p *plate) buildOutPath(filepath string) string {
	return path.Join(p.outPath, filepath)
}

func (p *plate) openTemplate(name string) (*template.Template, error) {
	t := template.New("")
	f, err := os.Open(p.buildTemplatePath(name))
	if err != nil {
		return t, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return t, err
	}

	return t.Parse(string(content))
}

func (p *plate) availableTemplates() []string {
	pattern := path.Join(p.srcPath, fmt.Sprintf("*%s", templatesExtension))
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	var names []string

	for _, path := range paths {
		name, err := filepath.Rel(p.srcPath, path)
		if err != nil {
			log.Fatal(err)
		}

		names = append(names, name[0:len(name)-len(templatesExtension)])
	}

	return names
}

func (p *plate) execute(name string) error {
	t, err := p.openTemplate(name)
	if err != nil {
		return err
	}

	for _, tpl := range t.Templates() {
		name := tpl.Name()
		if name != "" {
			path := p.buildOutPath(name)
			dir := filepath.Dir(path)
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}

			f, err := os.Create(path)
			if err != nil {
				return err
			}

			buf := bytes.NewBuffer([]byte{})
			tpl.Execute(buf, nil)
			tplContent := strings.TrimSpace(buf.String())
			io.WriteString(f, tplContent)
		}
	}

	return nil
}

func chooseTemplate(p *plate) string {
	templates := p.availableTemplates()
	fmt.Printf("Available templates:\n\n")
	for i, path := range templates {
		fmt.Printf("  %d - %v\n", i+1, path)
	}

	fmt.Printf("\nChoose your template [1-%d]: ", len(templates))

	var i int
	fmt.Scanf("%d", &i)

	if i < 1 || i > len(templates) {
		return chooseTemplate(p)
	}

	return templates[i-1]
}

func main() {
	var tplName string

	flag.StringVar(&tplName, "t", "", "template name")
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	templatesPath := path.Join(usr.HomeDir, templatesFolder)

	args := os.Args

	if len(args) != 2 {
		fmt.Printf("Usage:\n  %s PROJECT_PATH\n", args[0])
		os.Exit(1)
	}

	p := newPlate(templatesPath, args[1])
	name := chooseTemplate(p)
	p.execute(name)
}
