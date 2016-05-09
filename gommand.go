package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/imports"
)

func run(fileName string) (string, error) {
	out, err := exec.Command("go", "run", fileName).CombinedOutput()
	return string(out), err
}

func tempFile() (*os.File, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file, err := ioutil.TempFile(curDir, "temp")
	if err != nil {
		return nil, err
	}
	file.Close()

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

func usage() {
	fmt.Println("Usage: gommand [code]")
	fmt.Println("Example: gommand 'name := \"Sno6\"; fmt.Println(name)'")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	file, err := tempFile()
	if err != nil {
		log.Fatalf("main: error creating temp file: %v\n", err)
	}
	defer func() {
		file.Close()

		if err = os.Remove(file.Name()); err != nil {
			log.Printf("main: error removing temp file: %v\n", err)
		}
	}()

	code := os.Args[1]
	if code == "" {
		usage()
	}

	// bp holds the go boiler plate code with user inputted code added.
	bp := fmt.Sprintf("package main\nfunc main() {\n\t%v\n}", code)

	// Write go code to temp file and add missing imports.
	// Use Printf over Fatalf so removing of temp file will run through defer.
	if err = ioutil.WriteFile(file.Name(), []byte(bp), 0644); err != nil {
		log.Printf("main: error writing code to temp file: %v\n", err)
	}
	if err = editImports(file); err != nil {
		log.Printf("main: error editing imports: %v\n", err)
	}

	out, err := run(file.Name())
	if err != nil {
		log.Printf("main: error running go code query: %v\n", err)
	}
	if out == "" {
		return
	}
	fmt.Println(strings.TrimSpace(out))
}
