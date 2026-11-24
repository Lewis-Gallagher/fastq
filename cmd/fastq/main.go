package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Lewis-Gallagher/fastq/fastq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: fastq <file.fastq[.gz]>")
	}
	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func run(path string) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.Reader = f

	if strings.HasSuffix(path, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gz.Close()
		r = gz
	}

	p := fastq.NewParser(r)

	for {
		rec, err := p.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s\n%s\n\n", rec.ReadId, rec.Info)
	}
}
