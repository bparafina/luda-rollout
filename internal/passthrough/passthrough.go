package passthrough

import (
	"fmt"
	"os"
	"os/exec"
)

// Run looks up kubectl on PATH, builds a "kubectl rollout <args...>" command,
// and executes it with stdio connected directly. It preserves the exit code
// via *exec.ExitError.
func Run(args []string) error {
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return fmt.Errorf("kubectl not found on PATH: %w", err)
	}

	cmdArgs := append([]string{"rollout"}, args...)
	cmd := exec.Command(kubectlPath, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
