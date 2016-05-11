package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

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
	// Read the code from the temp go file.
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

	// Write the edited file into the original temp file.
	if err = ioutil.WriteFile(fileName, res, 0644); err != nil {
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	if len(os.Args) < 2 {
		usage()
	}
	code := os.Args[1]
	if code == "" {
		usage()
	}

	file, err := tempFile()
	go func() {
		for sig := range c {
			fmt.Println("Recieved sig", sig)
			os.Remove(file.Name())
		}
	}()

	if err != nil {
		log.Fatalf("main: error creating temp file: %v\n", err)
	}
	file.Close()

	defer func() {
		if err = os.Remove(file.Name()); err != nil {
			log.Printf("main: error removing temp file: %v\n", err)
		}
	}()

	// bp holds the go boiler plate code with user inputted code added.
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
