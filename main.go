package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const pythonBin string = "C:\\msys64\\usr\\bin\\python.exe"
const ansiblePlaybookBin string = "C:\\msys64\\usr\\bin\\ansible-playbook"

var opts struct {
	ExtraVars []string `short:"e" long:"extra-vars" description:"extra vars"`
	Inventory string   `short:"i" long:"inventory" description:"inventory"`
	Version   bool     `long:"version" description:"version"`
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	pythonArgs := []string{ansiblePlaybookBin, "-v"}

	for _, exVar := range opts.ExtraVars {
		if strings.HasPrefix(exVar, "ansible_ssh_private_key_file") {
			keyVal := strings.Split(exVar, "=")
			msysPath := toMsysPath(keyVal[1])
			pythonArgs = append(pythonArgs, "-e", fmt.Sprintf("'%s'", "ansible_ssh_private_key_file="+msysPath))
		} else {
			pythonArgs = append(pythonArgs, "-e", fmt.Sprintf("'%s'", exVar))
		}
	}

	if opts.Inventory != "" {
		pythonArgs = append(pythonArgs, "-i", opts.Inventory)
	}

	if len(args) > 0 {
		playbookPath := toRelPath(args[0])
		pythonArgs = append(pythonArgs, playbookPath)
	}

	if opts.Version {
		pythonArgs = append(pythonArgs, "--version")
	}

	fmt.Printf("args: %v\n", pythonArgs)
	command := exec.Command(pythonBin, pythonArgs...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Run()
	os.Exit(command.ProcessState.ExitCode())
}

func toMsysPath(winPath string) (msysPath string) {
	msysPath = strings.ReplaceAll(winPath, "\\", "/")
	msysPath = strings.ReplaceAll(msysPath, ":", "")
	msysPath = "/" + msysPath
	return
}

func toRelPath(winAbsPath string) (relPath string) {
	currentDir, _ := os.Getwd()
	relPath, _ = filepath.Rel(currentDir, winAbsPath)
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	return
}
