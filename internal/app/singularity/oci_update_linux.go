// Copyright (c) 2018-2021, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package singularity

import (
	"os"
	"os/exec"

	"github.com/sylabs/singularity/pkg/sylog"
)

// OciUpdate updates container cgroups resources
func OciUpdate(containerID string, args *OciArgs) error {
	runcArgs := []string{
		"--root", RuncStateDir,
		"update",
		"-r", args.FromFile,
		containerID,
	}

	cmd := exec.Command(runc, runcArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdout
	sylog.Debugf("Calling runc with args %v", runcArgs)
	return cmd.Run()
}
