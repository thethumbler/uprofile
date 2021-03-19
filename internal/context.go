package internal

import (
	"fmt"
	"os"
)

type Context struct {
	User        string
	ProfilesDir string
}

func GetContext() Context {
	user := os.Getenv("USER")

	if user == "" {
		fmt.Fprintln(os.Stderr, "Got empty $USER env")
		os.Exit(1)
	}

	profilesDir := fmt.Sprintf("/home/%s.profiles", user)

	if info, err := os.Stat(profilesDir); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Unable to access '%s', is context initialized?\n", profilesDir)
		os.Exit(1)
	}

	return Context{
		User:        user,
		ProfilesDir: profilesDir,
	}
}
