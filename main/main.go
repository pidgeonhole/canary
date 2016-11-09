package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"io"
	"log"
	"time"
	"github.com/yi-jiayu/esd-tutor-runner-python/lib"
)

var (
	app = kingpin.New("esd-tutor-runner-python", "Utility program to run untrusted, user-submitted python source code in a Docker container.")

	// the source file must be provided
	sfile = app.Arg("src-file", "Submitted source code").Required().File()

	// docker image to use
	img = app.Flag("image", "Docker image to use").Required().String()

	// input can be passed as a file or through stdin
	ifile = app.Flag("in-file", "Input file. Defaults to standard input.").Short('i').File()

	// output can be written to a file or from stdout
	ofile = app.Flag("out-file", "Output file. Defaults to standard output.").Short('o').OpenFile(os.O_CREATE | os.O_WRONLY, 666)

	// maximum time submission is allowed to run
	intTimeout = app.Flag("timeout", "Timeout in seconds").Short('t').Default("10").Int()

)

func main() {
	app.Version("0.1.0")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var in io.Reader
	var out io.Writer

	if *ifile == nil {
		in = os.Stdin
	} else {
		in = *ifile
	}

	if *ofile == nil {
		out = os.Stdout
	} else {
		out = *ofile
	}

	timeout := time.Duration(*intTimeout) * time.Second

	runner := lib.Runner{
		SourceFile: *sfile,
		In: in,
		Out: out,
		Timeout: timeout,
		Image: *img,
	}

	err := runner.Run()
	if err != nil {
		log.Println(err)
	}
}