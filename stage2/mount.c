// Copyright 2014-2015 Apcera Inc. All rights reserved.
//
// Portions of this file are based on code from:
//   https://github.com/arachsys/containers
//
// Code is licensed under MIT.
// Copyright 2013 Chris Webb <chris@arachsys.com>

#ifndef INITD_SPAWNER_MOUNT_C
#define INITD_SPAWNER_MOUNT_C

#define _GNU_SOURCE

#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>

#include <sys/mount.h>
#include <sys/syscall.h>

#include "spawner.h"

static char *root;

char *tmpdir(void) {
	char *dir;

	if (!(dir = strdup("/tmp/XXXXXX")))
		error(1, errno, "strdup");
	else if (!mkdtemp(dir))
		error(1, errno, "Failed to create temporary directory");
	return dir;
}

static void bindnode(char *src, char *dst) {
	int fd;

	if ((fd = open(dst, O_WRONLY | O_CREAT, 0600)) >= 0)
		close(fd);
	if (mount(src, dst, NULL, MS_BIND, NULL) < 0)
		error(1, errno, "Failed to bind %s into new %s filesystem", src, dst);
}

void createroot(char *src, bool privileged) {
	mode_t mask;
	pid_t child;
	int res;
	int console;

	mask = umask(0);

	// Create a temp directory that will contain the new root.
	root = tmpdir();

	// Mount the source to the root temp directory.
	if (mount(src, root, NULL, MS_BIND | MS_REC, NULL) < 0)
		error(1, errno, "Failed to bind new root filesystem");
	else if (chdir(root) < 0)
		error(1, errno, "Failed to enter new root filesystem");

	// Setup /dev as tmpfs mounts within the container
	mkdir("dev" , 0755);
	if (privileged) {
		if (mount("devtmpfs", "dev", "devtmpfs", 0, "") < 0)
			error(1, errno, "Failed to mount /dev devtmpfs in new root filesystem");
	} else {
		if (mount("tmpfs", "dev", "tmpfs", 0, "mode=0755") < 0)
			error(1, errno, "Failed to mount /dev tmpfs in new root filesystem");

		// Populate /dev within the container
		bindnode("/dev/full", "dev/full");
		bindnode("/dev/null", "dev/null");
		bindnode("/dev/random", "dev/random");
		bindnode("/dev/tty", "dev/tty");
		bindnode("/dev/urandom", "dev/urandom");
		bindnode("/dev/zero", "dev/zero");

		res = symlink("pts/ptmx", "dev/ptmx");
		res = symlink("/proc/kcore", "dev/core");
		res = symlink("/proc/self/fd", "dev/fd");
		res = symlink("console", "dev/kmsg");

		res = symlink("fd/0", "dev/stdin");
		res = symlink("fd/1", "dev/stdout");
		res = symlink("fd/2", "dev/stderr");
	}

	// setup /dev/pts and /dev/shm
	mkdir("dev/pts", 0755);
	if (mount("devpts", "dev/pts", "devpts", 0, "newinstance,ptmxmode=0666") < 0)
		error(1, errno, "Failed to mount /dev/pts in new root filesystem");
	mkdir("dev/shm", 0755);
	if (mount("tmpfs", "dev/shm", "tmpfs", 0, "mode=1777,size=65536k") < 0)
		error(1, errno, "Failed to mount /dev/pts in new root filesystem");

	// Setup /tmp within the container
	mkdir("tmp", 0777);
	if (mount("tmpfs", "tmp", "tmpfs", 0, "mode=0755") < 0)
		error(1, errno, "Failed to mount /tmp tmpfs in new root filesystem");
	umask(mask);

	// Setup /dev/console
	console = getconsole();
	bindnode(ptsname(console), "dev/console");
}

void enterroot(bool privileged) {
	if (chdir(root) < 0)
		error(1, errno, "Failed to chdir into the new root");
	if (mkdir("host", 0755) < 0)
		error(1, errno, "Failed to create host to place old filesystem at");
	if (syscall(__NR_pivot_root, ".", "host") < 0)
		error(1, errno, "Failed to pivot into new root filesystem");
	if (chdir("/") < 0 )
		error(1, errno, "Failed to detach old root filesystem");

	if (!privileged) {
		if (umount2("/host", MNT_DETACH) < 0)
			error(1, errno, "Failed to detach old root filesystem");
		rmdir("/host");
	}
}

void mountproc(void) {
	mode_t mask;

	mask = umask(0);
	mkdir("proc" , 0755);
	umask(mask);

	if (mount("proc", "proc", "proc", 0, NULL) < 0)
		error(1, errno, "Failed to mount /proc in new root filesystem: %s", strerror(errno));
}

#endif
