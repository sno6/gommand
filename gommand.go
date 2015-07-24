package main

import(
	"os"
	"os/exec"
	"io/ioutil"
	"fmt"
	"log"
	"strings"

	"golang.org/x/tools/imports"
)

func run(name string) (string, error) {
	out, err := exec.Command("go", "run", name).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), err
}

func tempFile() (*os.File, error) {
	curDir, err := os.Getwd()
	file, err := ioutil.TempFile(curDir, "temp")
	if err != nil {
		return nil, err
	}

	// Add the .go suffix to the temp file.
	if err = os.Rename(file.Name(), file.Name() + ".go"); err != nil {
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
	bp := fmt.Sprintf("package main\nfunc main() {\n\t%v\n}", code)

	if err = ioutil.WriteFile(file.Name(), []byte(bp), 0644); err != nil {
		log.Printf("main: error writing code to temp file: %v\n", err)
	}

	if err = editImports(file); err != nil {
		log.Printf("main: error editing imports: %v\n", err)
	}

	out, err := run(file.Name())
	if err != nil {
		fmt.Println("There was an error in your go code.")
	}

	if out != "" {
		fmt.Println(strings.TrimSpace(out))
	}
}
