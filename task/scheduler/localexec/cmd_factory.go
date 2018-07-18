package localexec

import "os/exec"

type cmdFactory func(...string) *exec.Cmd
