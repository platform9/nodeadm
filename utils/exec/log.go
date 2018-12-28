package exec

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	log "github.com/platform9/nodeadm/pkg/logrus"
)

// LogRun runs the command, streaming stdout and stderr to the log. Stdout and
// stderr lines may be interleaved.
func LogRun(cmd *exec.Cmd) error {
	op, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to open stdout pipe: %v", err)
	}
	ep, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to open stderr pipe: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		logAllWithPrefix("stdout", op)
	}()
	go func() {
		defer wg.Done()
		logAllWithPrefix("stderr", ep)
	}()
	cmd.Start()
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command %q failed: %v", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func logAllWithPrefix(prefix string, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		log.Infof("%v: %v", prefix, s.Text())
	}
}
