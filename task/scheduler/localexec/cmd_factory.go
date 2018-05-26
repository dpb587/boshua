package localexec

import "os/exec"

type CmdFactory func(...string) *exec.Cmd
