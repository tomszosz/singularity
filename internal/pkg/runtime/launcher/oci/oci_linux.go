// Copyright (c) 2018-2022, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

// Includes code from https://github.com/containers/podman
// Released under the Apache License Version 2.0

package oci

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/sylabs/singularity/internal/pkg/util/fs"
	"github.com/sylabs/singularity/internal/pkg/util/user"
	"github.com/sylabs/singularity/pkg/syfs"
	"github.com/sylabs/singularity/pkg/util/fs/lock"
)

const (
	// Absolute path for the runc state
	runcStateDir = "/run/singularity-oci"
	// Relative path inside ~/.singularity for conmon and singularity state
	ociPath = "oci"
	// State directory files
	containerPidFile = "container.pid"
	containerLogFile = "container.log"
	runcLogFile      = "runc.log"
	conmonPidFile    = "conmon.pid"
	bundleLink       = "bundle"
	// Files in the OCI bundle root
	bundleLock   = ".singularity-oci.lock"
	attachSocket = "attach"
	// Timeouts
	createTimeout = 30 * time.Second
)

// stateDir returns the path to container state handled by conmon/singularity
// (as opposed to runc's state in RuncStateDir)
func stateDir(containerID string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	u, err := user.CurrentOriginal()
	if err != nil {
		return "", err
	}

	configDir, err := syfs.ConfigDirForUsername(u.Name)
	if err != nil {
		return "", err
	}

	rootPath := filepath.Join(configDir, ociPath)
	containerPath := filepath.Join(hostname, containerID)
	path, err := securejoin.SecureJoin(rootPath, containerPath)
	if err != nil {
		return "", err
	}
	return path, err
}

// lockBundle creates a lock file in a bundle directory
func lockBundle(bundlePath string) error {
	bl := path.Join(bundlePath, bundleLock)
	_, err := os.Stat(bl)
	if err == nil {
		return fmt.Errorf("bundle is locked by another process")
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("while stat-ing lock file: %w", err)
	}

	fd, err := lock.Exclusive(bundlePath)
	if err != nil {
		return fmt.Errorf("while acquiring directory lock: %w", err)
	}
	defer lock.Release(fd)

	err = fs.EnsureFileWithPermission(bl, 0o600)
	if err != nil {
		return fmt.Errorf("while creating lock file: %w", err)
	}
	return nil
}

// releaseBundle removes a lock file in a bundle directory
func releaseBundle(bundlePath string) error {
	bl := path.Join(bundlePath, bundleLock)
	_, err := os.Stat(bl)
	if os.IsNotExist(err) {
		return fmt.Errorf("bundle is not locked")
	}
	if err != nil {
		return fmt.Errorf("while stat-ing lock file: %w", err)
	}

	fd, err := lock.Exclusive(bundlePath)
	if err != nil {
		return fmt.Errorf("while acquiring directory lock: %w", err)
	}
	defer lock.Release(fd)

	err = os.Remove(bl)
	if err != nil {
		return fmt.Errorf("while removing lock file: %w", err)
	}
	return nil
}