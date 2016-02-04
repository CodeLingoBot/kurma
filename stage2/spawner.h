// Copyright 2014-2016 Apcera Inc. All rights reserved.

#ifndef INITD_SPAWNER_SPAWNER_H
#define INITD_SPAWNER_SPAWNER_H

#include <errno.h>
#include <error.h>
#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <sys/capability.h>

// This is the private structure used within the clone_* calls. This contains
// a copy of all data used as well as the stack. This is allocated via a
// call to mmap.
typedef struct clone_destination_data {
	// the file that will be executed.
	char *command;

	// An array of strings passed as the argv for the command. The first
	// element of this list should always be populated. This array
	// is NULL terminated.
	char **args;

	// The environment to be added as key=value strings and NULL terminated.
	char **environment;

	// The list of cgroups tasks files that this task should join before
	// execing. This is an NULL terminated array like args or environment.
	char **tasksfiles;

	// The file descriptor that will be duplicated into the stdin position.
	int stdinfd;

	// The file descriptor that will be duplicated into the stdout position.
	int stdoutfd;

	// The file descriptor that will be duplicated into the stderr position.
	int stderrfd;

	// The setting for whether proc should be mounted for the new container.
	bool setup_proc;

	// The IPC namespace that should be joined after cloning.
	char *ipcnamespace;

	// The Mount namespace that should be joined after cloning.
	char *mountnamespace;

	// The network namespace that should be joined after cloning.
	char *networknamespace;

	// The pid namespace that should be joined after cloning.
	char *pidnamespace;

	// The UTS namespace that should be joined after cloning.
	char *utsnamespace;

	// The User namespace that should be joined after cloning.
	char *usernamespace;

	// Setup a new IPC namespace on clone.
	bool new_ipc_namespace;

	// Setup a new mount namespace on clone.
	bool new_mount_namespace;

	// Setup a new Network namespace on clone.
	bool new_network_namespace;

	// Setup a new pid namespace on clone.
	bool new_pid_namespace;

	// Setup a new UTS namespace on clone.
	bool new_uts_namespace;

	// Setup a new user namespace on clone.
	bool new_user_namespace;

	// Specifies whether the container should be setup in a privileged mode.
	bool privileged;

	// Tells the spawner to chroot into the directory.
	bool chroot;

	// The capabilities to apply to the process within the container.
	char *capabilities;

	// The UID mapping to write to the container's uid_map file
	char *uidmap;

	// The GID mapping to write to the container's gid_map file
	char *gidmap;

	// Maximum number of open files in the spawned process.
	int max_open_files;

	// Maximum numnber of processes that can be created by the spawned process.
	int max_processes;

	// The directory for the container's filesystem
	char *container_directory;

	// The directory where the container's filesystem should be bind mounted to
	char *bind_directory;

	// The user and to run the stage3 process as
	char *user;
	char *group;

	// True if this process should double fork in order to become a child of
	// spanwer rather than the calling process.
	bool detach;
} clone_destination_data;

// clone.c
void spawn_child(clone_destination_data *args);
static void setup_container(clone_destination_data *args, pid_t child);

// control.c
void dup_filedescriptors(int stdinfd, int stdoutfd, int stderrfd);
void closefds();
void joincgroups(char *tasksfiles[]);
void joinnamespace(char *filename);
int flags_for_clone(clone_destination_data *args);

// mount.c
char *tmpdir(void);
void privatize_namespace();
void createroot(char *src, char *dst, bool privileged);
void enterroot();
void mountproc(void);

// util.c
char *append(char **destination, const char *format, ...);
char *string(const char *format, ...);
void spawner_print_time(FILE *fd);
void writemap(pid_t pid, char *type, char *map);
void waitforstop(pid_t child);
void waitforexit(pid_t child);
int uidforuser(char *user);
int gidforgroup(char *group);
void dropBoundingCapabilities(cap_t newcap);

// -------
// Logging
// -------

// Set to 1 if debugging should be enabled.
extern bool spawner_debugging;

#define DEBUG(...)								\
	do {											\
	if (spawner_debugging) {						\
		fflush(NULL);								\
		spawner_print_time(stdout);						\
		fprintf(stdout, __VA_ARGS__);				\
		fflush(stdout);							\
	}											\
	} while(0)
#define INFO(...)								\
	do {											\
	spawner_print_time(stdout);							\
	fprintf(stdout, __VA_ARGS__);				\
	fflush(stdout);								\
	} while(0)
#define ERROR(...)								\
	do {											\
	spawner_print_time(stderr);							\
	fprintf(stderr, __VA_ARGS__);				\
	fflush(stderr);								\
	} while(0)
#define FATAL(...)								\
	do {											\
	spawner_print_time(stderr);							\
	fprintf(stderr, __VA_ARGS__);				\
	fflush(stderr);								\
	exit(1)										\
		} while(0)

#endif
