package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jsgilmore/mount"
)

func MountedProfiles(context *Context) []string {
	var mountedProfiles []string

	if mounts, err := mount.Mounts(); err != nil {
		log.Fatal(err)
	} else {
		for _, mount := range mounts {
			if mount.Filesystem == "fuse.fuse-overlayfs" {
				var match string
				n, _ := fmt.Sscanf(mount.Path, fmt.Sprintf("%s/%%s", context.ProfilesDir), &match)
				if n == 1 {
					profile := strings.Split(match, "/")[0]
					mountedProfiles = append(mountedProfiles, profile)
				}
			}
		}
	}

	return mountedProfiles
}

func IsMounted(context *Context, profile string) bool {
	mountedProfiles := MountedProfiles(context)

	for _, mountedProfile := range mountedProfiles {
		if mountedProfile == profile {
			return true
		}
	}

	return false
}

func MountProfile(context *Context, profile string) error {
	if IsMounted(context, profile) {
		fmt.Fprintf(os.Stderr, "profile %s is already mounted\n", profile)
		return nil
	}

	lowerDir := fmt.Sprintf("/home/%s", context.User)
	upperDir := fmt.Sprintf("%s/%s/upperdir", context.ProfilesDir, profile)
	workDir := fmt.Sprintf("%s/%s/workdir", context.ProfilesDir, profile)
	mergedDir := fmt.Sprintf("%s/%s/merged", context.ProfilesDir, profile)

	mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	return exec.Command("fuse-overlayfs", "-o", mountOptions, mergedDir).Run()
}

func UmountProfile(context *Context, profile string) error {
	if !IsMounted(context, profile) {
		fmt.Fprintf(os.Stderr, "profile %s is not mounted\n", profile)
		return nil
	}

	mountpoint := fmt.Sprintf("%s/%s/merged", context.ProfilesDir, profile)
	return exec.Command("umount", mountpoint).Run()
}
