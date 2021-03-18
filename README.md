ualias: manage multiple profiles (aliases) for a single linux user
===

Problem
-----
I work on mutliple contexts (mutliple projects, freelancing, etc). I need to manage multiple identites (mutliple SSH identites, k8s contexts, git config, etc) and I want to be able to jump between them quickly.

Solutions
---
### Use per-command flags to switch contexts
* Use `kubectx` for k8s
```
$ kubectx context1
$ kubectl ...
$ kubectx context2
$ kubectl ...
```
* Use `ssh -i <private key>` to use non-default private key
* Use per repo git configuration

### Use multiple users
You can create a user for each context and use the default/global configuration files in each user home directory. This works well for completely isolated work. However I often want to share files between all users (like my VIM configuration), which is hard to maintain.

### Use ualias
ualias allows you to have multiple identites on the same user, all identites have their own home directories that is layered on top of the user home directory, so all files that don't need to be edited are shared.

```
$ ualias create context1
$ ualias mount context1
$ ualias jump context1
(context1) $ git config --global user.email example@example.com # doesn't override original user global config
```

Principle
---
ualias uses `fuse-overlayfs` to mount identity home directory (similar to docker but without full isolation) and overrides `HOME` environment variable. Any writes in the identity home directory are not committed to original user home directory, keeping them isolated.

ualias also uses `unshare` to jump to the identity context. While not immediately necessary, using `unshare` allows us to add more complex features later on.


Building
---
Make sure that `fuse-overlayfs` and `unshare` are installed.

```
$ go build
$ ./ualias ...
```

It's prefered to move `ualias` binary to `~/.local/bin/` and have `~/.local/bin/` in your `PATH` environment variable.


Usage
---
We need a directory `/home/$USER.aliases` to host all identites that is accessable by `$USER`.

```
$ sudo mkdir /home/$USER.aliases
$ sudo chown $USER: /home/$USER.aliases
```

Then we can use `ualias` to create, mount and jump to an alias
```
$ ualias create context1
$ ualias mount context1
$ ualias jump context1
(context1) $ ...
```

an alias can be unmounted later with
```
$ ualias umount context1
```

if the alias is no longer needed, we can delete it (this will delete all the modifications in alias home directory though, so make sure to backup anything you need before that)

```
$ ualias delete context1 # context1 must be unmounted first
```