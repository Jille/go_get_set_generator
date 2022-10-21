package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jille/go_get_set_generator/get_set_generate/internal"
)

var (
	noSetters = flag.Bool("no-setters", false, "Whether to not generate setters")
	noGetters = flag.Bool("no-getters", false, "Whether to not generate getters")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("No source files given")
	}

	for _, sf := range flag.Args() {
		if err := generate(sf); err != nil {
			log.Fatalf("Error handling %s: %v", sf, err)
		}
	}
}

func generate(sourceFile string) error {
	fileInfo, err := internal.ParseFile(sourceFile)
	if err != nil {
		return fmt.Errorf("Loading input failed: %v", err)
	}

	fileName, _ := filepath.Abs(sourceFile)
	f, err := os.Create(strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "_getter_setter.go")
	if err != nil {
		return fmt.Errorf("Failed opening destination file: %v", err)
	}

	g := internal.Generator{
		Setters: !*noSetters,
		Getters: !*noGetters,
	}
	if err := g.Generate(fileInfo); err != nil {
		return fmt.Errorf("Failed generating getter setter: %v", err)
	}

	if _, err := f.Write(g.Output()); err != nil {
		return fmt.Errorf("Failed writing to destination: %v", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close destination file: %v", err)
	}
	return nil
}
