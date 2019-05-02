package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"golang.org/x/tools/imports"
)

func run(fileName string) error {
	cmd := exec.Command("go", "run", fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

func editImports(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	// res holds the go file re-written with imports added.
	opt := &imports.Options{}
	res, err := imports.Process(fileName, data, opt)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, res, 0644)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: gommand [code]")
	os.Exit(2)
}

func clean(f *os.File) {
	f.Close()
	if err := os.Remove(f.Name()); err != nil {
		log.Printf("main: error removing tempfile: %v\n", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	code := os.Args[1]
	if code == "" {
		usage()
	}

	file, err := tempFile()
	if err != nil {
		log.Fatalf("main: error creating temp file: %v\n", err)
	}
	defer clean(file)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		clean(file)
		os.Exit(1)
	}()

	// bp holds the go boiler plate code with user code added.
	bp := fmt.Sprintf("package main\nfunc main() {\n\t%v\n}", code)

	// Write go code to temp file and add missing imports.
	// Use Printf over Fatalf so removing of temp file will run through defer.
	if err = ioutil.WriteFile(file.Name(), []byte(bp), 0644); err != nil {
		log.Printf("main: error writing code to temp file: %v\n", err)
	}
	if err = editImports(file.Name()); err != nil {
		log.Printf("main: error editing imports: %v\n", err)
	}
	if err = run(file.Name()); err != nil {
		log.Printf("main: error running go code query: %v\n", err)
	}
}
