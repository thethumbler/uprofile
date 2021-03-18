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
	User       string
	AliasesDir string
}

func (ctx *Context) Create() {
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.Parse(os.Args[2:])

	if len(createCommand.Args()) != 1 {
		createCommand.Usage()
		return
	}

	alias := createCommand.Arg(0)

	dirsList := []string{
		fmt.Sprintf("%s/%s/upperdir", ctx.AliasesDir, alias),
		fmt.Sprintf("%s/%s/workdir", ctx.AliasesDir, alias),
		fmt.Sprintf("%s/%s/merged", ctx.AliasesDir, alias),
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

	alias := deleteCommand.Arg(0)
	aliasDir := fmt.Sprintf("%s/%s", ctx.AliasesDir, alias)
	os.RemoveAll(aliasDir)
}

func (ctx *Context) List() {
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCommand.Bool("a", false, "list all aliases (including not mounted)")
	listCommand.Parse(os.Args[2:])

	if *listAll {
		aliasesList, _ := ioutil.ReadDir(ctx.AliasesDir)
		for _, entry := range aliasesList {
			fmt.Println(entry.Name())
		}
	} else {
		mountedAliases, _ := mount.Mounts()
		for _, mount := range mountedAliases {
			if mount.Filesystem == "fuse.fuse-overlayfs" {
				var match string
				n, _ := fmt.Sscanf(mount.Path, fmt.Sprintf("%s/%%s", ctx.AliasesDir), &match)
				if n == 1 {
					alias := strings.Split(match, "/")[0]
					fmt.Println(alias)
				}
			}
		}
	}
}

func (ctx *Context) Mount() {
	mountCommand := flag.NewFlagSet("mount", flag.ExitOnError)
	mountCommand.Parse(os.Args[2:])

	alias := mountCommand.Arg(0)

	lowerDir := fmt.Sprintf("/home/%s", ctx.User)
	upperDir := fmt.Sprintf("%s/%s/upperdir", ctx.AliasesDir, alias)
	workDir := fmt.Sprintf("%s/%s/workdir", ctx.AliasesDir, alias)
	mergedDir := fmt.Sprintf("%s/%s/merged", ctx.AliasesDir, alias)

	mountOptions := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	exec.Command("fuse-overlayfs", "-o", mountOptions, mergedDir).Run()
}

func (ctx *Context) Umount() {
	umountCommand := flag.NewFlagSet("umount", flag.ExitOnError)
	umountCommand.Parse(os.Args[2:])

	alias := umountCommand.Arg(0)
	mountpoint := fmt.Sprintf("%s/%s/merged", ctx.AliasesDir, alias)
	exec.Command("umount", mountpoint).Run()
}

func (ctx *Context) Jump() {
	jumpCommand := flag.NewFlagSet("jump", flag.ExitOnError)
	jumpCommand.Parse(os.Args[2:])

	alias := jumpCommand.Arg(0)
	aliasDir := fmt.Sprintf("%s/%s", ctx.AliasesDir, alias)
	homeDir := fmt.Sprintf("%s/merged", aliasDir)

	if info, err := os.Stat(aliasDir); err != nil {
		if err.Error() == "no such file or directory" {
			fmt.Fprintf(os.Stderr, "no such alias: %s\n", alias)
		} else {
			fmt.Fprintf(os.Stderr, "unable to access alias '%s', '%s': %s\n", alias, aliasDir, err.Error())
		}

		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "alias '%s' is not created correctly, '%s' is not a directory.\n", alias, aliasDir)
		os.Exit(1)
	}

	env := append(os.Environ(),
		fmt.Sprintf("HOME=%s", homeDir),
		fmt.Sprintf("PS1=[%s/\\u@\\h \\W]\\$ ", alias),
	)

	binaryPath := "/usr/bin/unshare"
	if err := syscall.Exec(binaryPath, []string{"unshare", "-w", homeDir}, env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run 'unshare': %s", err.Error())
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: ualias <COMMAND> [FLAGS] [ARGS]")
	fmt.Println("\nCOMMAND:")
	fmt.Println("  create      create a new alias")
	fmt.Println("  delete      delete an existing alias")
	fmt.Println("  list        list aliases")
	fmt.Println("  mount       mount alias")
	fmt.Println("  umount      unmount a mounted alias")
	fmt.Println("  jump        jump to alias context")
}

func main() {
	user := os.Getenv("USER")

	if user == "" {
		fmt.Fprintln(os.Stderr, "Got empty $USER env")
		os.Exit(1)
	}

	context := Context{
		User:       user,
		AliasesDir: fmt.Sprintf("/home/%s.aliases", user),
	}

	if info, err := os.Stat(context.AliasesDir); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Unable to access '%s', is context initialized?\n", context.AliasesDir)
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
