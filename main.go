package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	input_file       string
	output_file      string
	verbouse         bool
	help             bool
	workers          int
	timeout_respones int
)
var input_data string
var output_data io.Writer

func init() {
	flag.BoolVar(&verbouse, "v", false, "Output fool log to StdOut (shorthand)")
	flag.BoolVar(&verbouse, "verbouse", false, "Output fool log to StdOut")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.StringVar(&input_file, "i", "", "Input web src for scraping data. If the flag is absent then input should from last argument.")
	flag.StringVar(&output_file, "o", "", "File for result output. If the flag is absent then output will to the StdOut.")
	flag.IntVar(&workers, "w", 5, "The number of workers working in the same time.")
	flag.IntVar(&timeout_respones, "t", 5, "The timeout in seconds for waiting a responses from web sites.")

}

func showHelp() {
	flag.VisitAll(func(f *flag.Flag) {
		if f.DefValue == "" {
			fmt.Printf("\t-%s: %s\n", f.Name, f.Usage)
		} else {
			fmt.Printf("\t-%s: %s (Default: %s)\n", f.Name, f.Usage, f.DefValue)
		}
	})
	os.Exit(0)
}

func parsFlags() {
	if help {
		showHelp()
	}
	if !verbouse {
		output_log, err := os.OpenFile("gscrape.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal("Cant open ouput file for loging gscrape.log.\n", err)
		}
		log.SetOutput(output_log)
	}

	if output_file == "" {
		output_data = os.Stdout
	} else {
		out, err := os.OpenFile(output_file, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal("Cant open ouput file", output_file, "\r\n", err)
		}
		output_data = out
	}

	switch {
	case input_file != "":
		data, err := os.ReadFile(input_file)
		if err != nil {
			panic(err)
		}
		input_data = string(data)
	case flag.Arg(0) != "":
		input_data = flag.Arg(0)
	default:
		fmt.Println("A flag is absent. The input file flag is expected...")
		showHelp()
	}
}

func run(HREF string, f *io.Writer) error {
	fout := bufio.NewWriter(*f)
	defer fout.Flush()
	pool := Newpool(workers)
	mutex := new(sync.Mutex)
	href_in := bufio.NewReader(strings.NewReader(HREF))

	for task := range NewTasks(href_in) {
		pool.worker(task.parse, NewOrgJson(task.num), fout, mutex)
	}

	pool.Wait_g.Wait()
	return pool.err
}

func main() {
	time_start := time.Now()
	flag.Parse()

	parsFlags()

	out := output_file
	if output_file == "" {
		out = "StdOut"
	}
	if input_file == "" {
		fmt.Printf("output_file: %s\nworkers: %d\n", out, workers)
	} else {
		fmt.Printf("output_file: %s\ninput_file: %s\nworkers: %d\n", out, input_file, workers)
	}

	run(input_data, &output_data)
	if true {
		fmt.Printf("END. Errors in log... Time %v ms.", time.Since(time_start).Milliseconds())
	} else {
		fmt.Printf("END. OK... Time %v ms.", time.Since(time_start).Milliseconds())
	}
}
