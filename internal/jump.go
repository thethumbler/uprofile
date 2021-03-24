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
		fmt.Sprintf("PS1=[\\u@%s \\W]\\$ ", profile),
	)

	binaryPath := "/usr/bin/unshare"
	if err := syscall.Exec(binaryPath, []string{"unshare", "-w", homeDir}, env); err != nil {
		return err
	}

	return nil
}
