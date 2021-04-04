#define _GNU_SOURCE
#include <dlfcn.h>
#include <pwd.h>
#include <unistd.h>
#include <stdlib.h>

struct passwd *getpwuid(uid_t uid) {
    struct passwd *(*libc_getpwuid)(uid_t uid) = dlsym(RTLD_NEXT, "getpwuid");
    struct passwd *passwd = libc_getpwuid(uid);

    if (uid == getuid()) {
        passwd->pw_dir = getenv("HOME");
    }

    return passwd;
}
