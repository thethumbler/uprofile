package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/jsgilmore/mount"
)

type Context struct {
	User        string
	ProfilesDir string
}

func (ctx *Context) Create() {
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.Parse(os.Args[2:])

	if len(createCommand.Args()) != 1 {
		createCommand.Usage()
		return
	}

	profile := createCommand.Arg(0)

	dirsList := []string{
		fmt.Sprintf("%s/%s/upperdir", ctx.ProfilesDir, profile),
		fmt.Sprintf("%s/%s/workdir", ctx.ProfilesDir, profile),
		fmt.Sprintf("%s/%s/merged", ctx.ProfilesDir, profile),
	}

	for _, path := range dirsList {
		os.MkdirAll(path, 0700)
	}
}

func (ctx *Context) Delete() {
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCommand.Parse(os.Args[2:])

	if len(deleteCommand.Args()) != 1 {
		deleteCommand.Usage()
		return
	}

	profile := deleteCommand.Arg(0)
	profileDir := fmt.Sprintf("%s/%s", ctx.ProfilesDir, profile)
	os.RemoveAll(profileDir)
}

func (ctx *Context) List() {
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCommand.Bool("a", false, "list all porfiles (including not mounted)")
	listCommand.Parse(os.Args[2:])

	if *listAll {
		profilesList, _ := ioutil.ReadDir(ctx.ProfilesDir)
		for _, entry := range profilesList {
			fmt.Println(entry.Name())
		}
	} else {
		mountedProfiles, _ := mount.Mounts()
		for _, mount := range mountedProfiles {
			if mount.Filesystem == "fuse.fuse-overlayfs" {
				var match string
				n, _ := fmt.Sscanf(mount.Path, fmt.Sprintf("%s/%%s", ctx.ProfilesDir), &match)
				if n == 1 {
					profile := strings.Split(match, "/")[0]
					fmt.Println(profile)
				}
			}
		}
	}
}

func (ctx *Context) Mount() {
	mountCommand := flag.NewFlagSet("mount", flag.ExitOnError)
	mountCommand.Parse(os.Args[2:])

	profile := mountCommand.Arg(0)

	lowerDir := fmt.Sprintf("/home/%s", ctx.User)
	upperDir := fmt.Sprintf("%s/%s/upperdir", ctx.ProfilesDir, profile)
	workDir := fmt.Sprintf("%s/%s/workdir", ctx.ProfilesDir, profile)
	mergedDir := fmt.Sprintf("%s/%s/merged", ctx.ProfilesDir, profile)

	mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	exec.Command("fuse-overlayfs", "-o", mountOptions, mergedDir).Run()
}

func (ctx *Context) Umount() {
	umountCommand := flag.NewFlagSet("umount", flag.ExitOnError)
	umountCommand.Parse(os.Args[2:])

	profile := umountCommand.Arg(0)
	mountpoint := fmt.Sprintf("%s/%s/merged", ctx.ProfilesDir, profile)
	exec.Command("umount", mountpoint).Run()
}

func (ctx *Context) Jump() {
	jumpCommand := flag.NewFlagSet("jump", flag.ExitOnError)
	jumpCommand.Parse(os.Args[2:])

	profile := jumpCommand.Arg(0)
	profileDir := fmt.Sprintf("%s/%s", ctx.ProfilesDir, profile)
	homeDir := fmt.Sprintf("%s/merged", profileDir)

	if info, err := os.Stat(profileDir); err != nil {
		if err.Error() == "no such file or directory" {
			fmt.Fprintf(os.Stderr, "no such profile: %s\n", profile)
		} else {
			fmt.Fprintf(os.Stderr, "unable to access profile '%s', '%s': %s\n", profile, profileDir, err.Error())
		}

		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "profile '%s' is not created correctly, '%s' is not a directory.\n", profile, profileDir)
		os.Exit(1)
	}

	env := append(os.Environ(),
		fmt.Sprintf("HOME=%s", homeDir),
		fmt.Sprintf("PS1=[%s/\\u@\\h \\W]\\$ ", profile),
	)

	binaryPath := "/usr/bin/unshare"
	if err := syscall.Exec(binaryPath, []string{"unshare", "-w", homeDir}, env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run 'unshare': %s", err.Error())
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: uprofile <COMMAND> [FLAGS] [ARGS]")
	fmt.Println("\nCOMMAND:")
	fmt.Println("  create      create a new profile")
	fmt.Println("  delete      delete an existing profile")
	fmt.Println("  list        list profiles")
	fmt.Println("  mount       mount profile")
	fmt.Println("  umount      unmount a mounted profile")
	fmt.Println("  jump        jump to profile context")
}

func main() {
	user := os.Getenv("USER")

	if user == "" {
		fmt.Fprintln(os.Stderr, "Got empty $USER env")
		os.Exit(1)
	}

	context := Context{
		User:        user,
		ProfilesDir: fmt.Sprintf("/home/%s.profiles", user),
	}

	if info, err := os.Stat(context.ProfilesDir); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Unable to access '%s', is context initialized?\n", context.ProfilesDir)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		context.Create()
	case "delete":
		context.Delete()
	case "list":
		context.List()
	case "mount":
		context.Mount()
	case "umount":
		context.Umount()
	case "jump":
		context.Jump()
	default:
		usage()
		os.Exit(1)
	}
}
