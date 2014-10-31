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

func main() {
	var tplName string

	flag.StringVar(&tplName, "t", "", "template name")
	flag.Parse()

	fmt.Printf(" --- %v\n", tplName)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	templatesPath := path.Join(usr.HomeDir, templatesFolder)

	args := os.Args

	fmt.Printf("%v", args)

	p := newPlate(templatesPath, args[1])
	err = p.execute(tplName)
	if err != nil {
		log.Fatal(err)
	}
}
