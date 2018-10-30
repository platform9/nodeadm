package deprecated

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	log "github.com/platform9/nodeadm/pkg/logrus"
)

func Run(rootDir string, cmdStr string, arg ...string) {
	if len(rootDir) > 0 {
		currentPath := os.Getenv("PATH")
		os.Setenv("PATH", currentPath+":"+rootDir)
		log.Infof("Updated PATH variable = %s", os.Getenv("PATH"))
		log.Infof("Running command %s %v", cmdStr, arg)
	}
	cmd := exec.Command(cmdStr, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("\nFailed to run command %s with error %v", cmdStr, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("\nFailed to get output of command %s with error %v", cmdStr, err)
	}
}

func PipeCmd(rootDir string, inputStr string, cmdStr string, arg ...string) {
	if len(rootDir) > 0 {
		currentPath := os.Getenv("PATH")
		os.Setenv("PATH", currentPath+":"+rootDir)
		log.Infof("Updated PATH variable = %s", os.Getenv("PATH"))
		log.Infof("Running command %s %v", cmdStr, arg)
	}

	echoCmd := exec.Command("echo", inputStr)
	kubectlCmd := exec.Command(cmdStr, arg...)

	reader, writer := io.Pipe()
	var buf bytes.Buffer

	echoCmd.Stdout = writer

	kubectlCmd.Stdin = reader
	kubectlCmd.Stderr = os.Stderr
	kubectlCmd.Stdout = &buf

	echoCmd.Start()
	kubectlCmd.Start()

	echoCmd.Wait()
	writer.Close()

	kubectlCmd.Wait()
	reader.Close()
}

func RunBestEffort(rootDir string, cmdStr string, arg ...string) {
	if len(rootDir) > 0 {
		currentPath := os.Getenv("PATH")
		os.Setenv("PATH", currentPath+":"+rootDir)
		log.Infof("Updated PATH variable = %s", os.Getenv("PATH"))
		log.Infof("Running command %s %v", cmdStr, arg)
	}
	cmd := exec.Command(cmdStr, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
}
