package lib

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Runner struct {
	SourceFile *os.File
	In         io.Reader
	Out        io.Writer
	Timeout    time.Duration
	Image      string
}

//func NewRunner(srcFile *os.File, in io.Reader, out io.Writer, timeout time.Duration, image string) (Runner, error) {
//
//}

func (r Runner) Run() error {
	sfileAbs, err := filepath.Abs(r.SourceFile.Name())
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path to source file")
	}

	bin := "docker"
	// -i              attach stdin to container
	// --rm            automatically remove container after it exits
	// --cap-drop=ALL  drop all capabilities from container
	// --net=none      take container offline
	// -v              mount source file inside container at /src/subm.py
	argstr := fmt.Sprintf("run -i --rm --cap-drop=ALL --net=none -v %s:/src/subm.py %s python /src/subm.py", sfileAbs, r.Image)
	args := strings.Split(argstr, " ")

	cmd := exec.Command(bin, args...)
	timeout := r.Timeout

	cmd.Stdout = r.Out
	cmd.Stdin = r.In
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to run docker container")
	}

	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
	})

	if err := cmd.Wait(); err != nil {
		timedOut := !timer.Stop()
		if timedOut {
			return errors.New("time limit exceeded")
		}

		return errors.Wrap(err, "error while running submission")
	}

	return nil
}
