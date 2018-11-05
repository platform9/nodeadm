package deprecated

import (
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
