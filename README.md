uprofile: manage multiple profiles for a single linux user
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

### Use uprofile
uprofile allows you to have multiple identites on the same user, all identites have their own home directories that is layered on top of the user home directory, so all files that don't need to be edited are shared.

```
$ uprofile create context1
$ uprofile mount context1
$ uprofile jump context1
(context1) $ git config --global user.email example@example.com # doesn't override original user global config
```

Principle
---
uprofile uses `fuse-overlayfs` to mount identity home directory (similar to docker but without full isolation) and overrides `HOME` environment variable. Any writes in the identity home directory are not committed to original user home directory, keeping them isolated.

uprofile also uses `unshare` to jump to the identity context. While not immediately necessary, using `unshare` allows us to add more complex features later on.


Building
---
Make sure that `fuse-overlayfs` and `unshare` are installed.

```
$ go build
$ ./uprofile ...
```

It's prefered to move `uprofile` binary to `~/.local/bin/` and have `~/.local/bin/` in your `PATH` environment variable.


Usage
---
We need a directory `/home/$USER.profiles` to host all identites that is accessable by `$USER`.

```
$ sudo mkdir /home/$USER.profiles
$ sudo chown $USER: /home/$USER.profiles
```

Then we can use `uprofile` to create, mount and jump to a profile
```
$ uprofile create context1
$ uprofile mount context1
$ uprofile jump context1
(context1) $ ...
```

a profile can be unmounted later with
```
$ uprofile umount context1
```

if the profile is no longer needed, we can delete it (this will delete all the modifications in profile home directory though, so make sure to backup anything you need before that)

```
$ uprofile delete context1 # context1 must be unmounted first
```

Advanced
---
I use monkey patching to overcome some of the applications that don't use/respect HOME env variable. For example, `OpenSSH` uses `getpwuid` to retrieve user home directory which we can't easily override. The source for monkey patching is located under `patch/` but I won't provide steps on how to use it to ensure that only users who really know what they are doing can use it :)
