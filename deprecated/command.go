package deprecated

import (
	"log"
	"os"
	"os/exec"
)

func Run(rootDir string, cmdStr string, arg ...string) {
	if len(rootDir) > 0 {
		currentPath := os.Getenv("PATH")
		os.Setenv("PATH", currentPath+":"+rootDir)
		log.Printf("Updated PATH variable = %s", os.Getenv("PATH"))
		log.Printf("Running command %s %v", cmdStr, arg)
	}
	cmd := exec.Command(cmdStr, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to run command %s with error %v\n", cmdStr, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Failed to get output of command %s with error %v\n", cmdStr, err)
	}
}

func RunBestEffort(rootDir string, cmdStr string, arg ...string) {
	if len(rootDir) > 0 {
		currentPath := os.Getenv("PATH")
		os.Setenv("PATH", currentPath+":"+rootDir)
		log.Printf("Updated PATH variable = %s", os.Getenv("PATH"))
		log.Printf("Running command %s %v", cmdStr, arg)
	}
	cmd := exec.Command(cmdStr, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
}
