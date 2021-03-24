package internal

import (
	"fmt"
	"os"
	"syscall"
)

func Jump(context *Context, profile string) error {
	profileDir := fmt.Sprintf("%s/%s", context.ProfilesDir, profile)
	homeDir := fmt.Sprintf("%s/merged", profileDir)

	if info, err := os.Stat(profileDir); err != nil {
		if err.Error() == "no such file or directory" {
			return fmt.Errorf("no such profile: %s", profile)
		} else {
			return fmt.Errorf("unable to access profile '%s', '%s': %s", profile, profileDir, err.Error())
		}
	} else if !info.IsDir() {
		return fmt.Errorf("profile '%s' is not created correctly, '%s' is not a directory", profile, profileDir)
	}

	if !IsMounted(context, profile) {
		if err := MountProfile(context, profile); err != nil {
			return err
		}
	}

	env := append(os.Environ(),
		fmt.Sprintf("HOME=%s", homeDir),
		fmt.Sprintf("UPROFILE=%s", profile),
	)

	if err := syscall.Unshare(0); err != nil {
		return err
	}

	if err := syscall.Chdir(homeDir); err != nil {
		return err
	}

	defaultShell := os.Getenv("SHELL")

	if defaultShell == "" {
		defaultShell = "/bin/sh"
	}

	if err := syscall.Exec(defaultShell, []string{defaultShell}, env); err != nil {
		return err
	}

	return nil
}
