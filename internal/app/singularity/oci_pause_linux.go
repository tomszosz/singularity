// Copyright (c) 2018-2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package singularity

import (
	"os"
	"os/exec"

	"github.com/sylabs/singularity/pkg/sylog"
)

// OciPause pauses processes in a container
func OciPause(containerID string) error {
	runcArgs := []string{
		"--root", RuncStateDir,
		"pause",
		containerID,
	}

	cmd := exec.Command(runc, runcArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdout
	sylog.Debugf("Calling runc with args %v", runcArgs)
	return cmd.Run()
}

// OciResume pauses processes in a container
func OciResume(containerID string) error {
	runcArgs := []string{
		"--root", RuncStateDir,
		"resume",
		containerID,
	}

	cmd := exec.Command(runc, runcArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdout
	sylog.Debugf("Calling runc with args %v", runcArgs)
	return cmd.Run()
}
