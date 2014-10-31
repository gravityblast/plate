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
	templatesFolder    = ".plates"
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

func (p *plate) setup() {
	os.MkdirAll(p.srcPath, 0777)
}

func (p *plate) buildTemplatePath(name string) string {
	filename := fmt.Sprintf("%s%s", name, templatesExtension)
	return path.Join(p.srcPath, filename)
}

func (p *plate) buildOutPath(filepath string) string {
	return path.Join(p.outPath, filepath)
}

func (p *plate) templateFuncs(args ...string) template.FuncMap {
	return template.FuncMap{
		"args": func(i int) string {
			if i >= len(args) {
				fmt.Printf("The current template requires Args[%d].\n", i)
				fmt.Printf("Current Args are:\n")
				for index, arg := range args {
					fmt.Printf("  %d: %s\n", index, arg)
				}
				os.Exit(1)
			}

			return args[i]
		},
	}
}

func (p *plate) openTemplate(name string, args ...string) (*template.Template, error) {
	t := template.New("")
	t.Funcs(p.templateFuncs(args...))

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

func (p *plate) execute(name string, args ...string) error {
	t, err := p.openTemplate(name, args...)
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
			err = tpl.Execute(buf, nil)
			if err != nil {
				return err
			}

			tplContent := strings.TrimSpace(buf.String())
			io.WriteString(f, tplContent)
		}
	}

	return nil
}

func chooseTemplate(p *plate) string {
	templates := p.availableTemplates()
	if len(templates) < 1 {
		log.Fatalf("No templates available in %s", p.srcPath)
	}
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

	if len(args) < 2 {
		fmt.Printf("Usage:\n  %s PROJECT_PATH\n", args[0])
		os.Exit(1)
	}

	p := newPlate(templatesPath, args[1])
	p.setup()
	name := chooseTemplate(p)
	err = p.execute(name, args...)
	if err != nil {
		log.Fatal(err)
	}
}
