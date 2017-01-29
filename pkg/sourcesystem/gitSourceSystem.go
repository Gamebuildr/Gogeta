package sourcesystem

import (
	"fmt"
	"os/exec"
)

func gitAddSource() {
	cmd := exec.Command("git", "clone", "--depth", "2", "testRepo", "testPath")

	commandErr := cmd.Start()
	if commandErr != nil {
		fmt.Printf("%d error: ", commandErr.Error())
	}
	fmt.Printf("clone successful")
}
