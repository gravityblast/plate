# plate

> Plate lets you quickly setup project files starting from templates defined in your `~/.plates` folder.

## Installation

```
go get github.com/pilu/plate
```

## Usage

Create a template for a Go hello world app:

`~/.plates/go-hello-world.plate`:
```
{{define "main.go"}}
package main

import "fmt"

func main() {
	fmt.Println("{{ ask "greeting"}}")
}
{{end}}

{{define "main_test.go"}}
package main

import "testing"

func TestFoo(t *testing.T) {
}
{{end}}
```

Run `plate .` and choose the go-hello-world plate.

Plate will create two files, `main.go` and `main_test.go`.

Before creating the `main.go` file, it will ask for a value for the variable `greeting`.

## Author

* [Andrea Franz](http://gravityblast.com)
