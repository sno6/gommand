package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"go/build"

	"golang.org/x/tools/imports"
)

const code_tmpl = `
package main

%s

func p(args ...interface{}) {
	fmt.Println(args...)
}

%s

func main() {
	%v
}
`

func run(name string, args []string) (string, error) {
	go_args := []string{"run", name}
	go_args = append(go_args, args...)
	out, err := exec.Command("go", go_args...).CombinedOutput()
	return string(out), err
}

func tempFile() (*os.File, error) {
	curDir, err := os.Getwd()
	file, err := ioutil.TempFile(curDir, "temp")
	if err != nil {
		return nil, err
	}

	// Add the .go suffix to the temp file.
	if err = os.Rename(file.Name(), file.Name()+".go"); err != nil {
		return nil, err
	}
	return os.Open(file.Name() + ".go")
}

func editImports(file *os.File) error {
	// Read the code from the temp go file.
	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	// res holds the go file re-written with imports added.
	opt := &imports.Options{}
	res, err := imports.Process(file.Name(), data, opt)
	if err != nil {
		return err
	}

	// Write the edited file into the original temp file.
	if err = ioutil.WriteFile(file.Name(), res, 0644); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gommand 'name := \"Sno\"; fmt.Println(name)'")
		return
	}

	file, err := tempFile()
	if err != nil {
		log.Printf("main: error creating temp file: %v\n", err)
	}
	defer file.Close()

	defer func() {
		if err = os.Remove(file.Name()); err != nil {
			log.Printf("main: error removing temp file: %v\n", err)
		}
	}()

	code := os.Args[1]
	// check github.com/k0kubun/pp is installed
	pp_import, pp_func := "", ""
	_, err = build.Import("github.com/k0kubun/pp", "", build.FindOnly)
	if err == nil {
		pp_import = `import (pp_dumper "github.com/k0kubun/pp")`
		pp_func = `
		func pp(args ...interface{}) {
			pp_dumper.Print(args...)
		}
		`
	}

	// bp holds the go boiler plate code and the added user input.
	bp := fmt.Sprintf(code_tmpl, pp_import, pp_func, code)

	if err = ioutil.WriteFile(file.Name(), []byte(bp), 0644); err != nil {
		log.Printf("main: error writing code to temp file: %v\n", err)
	}

	if err = editImports(file); err != nil {
		fmt.Printf("main: error editing imports: %v\n", err)
		return
	}

	out, err := run(file.Name(), os.Args[2:])
	if err != nil {
		fmt.Println(out)
		return
	}
	fmt.Println(strings.TrimSpace(out))
}
