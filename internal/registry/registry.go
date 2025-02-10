package registry

import (
	"os"
	"os/exec"
)

// Executes docker push
func PushImageToRegistry(image string) error {
	pushCmd := exec.Command("docker", "push", image)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	return pushCmd.Run()
}
