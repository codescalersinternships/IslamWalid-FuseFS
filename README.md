# Fuse Filesystem

Fuse filesystem is a userspace file system, which lets the non-previliged users to create their own filesystem without editing the kernel code.
This is done by running the actual filesystem code in the userspace while the FUSE interface works as a bridge between the userspace and the kernel.

## Functionality
- `NewFS(userStructReference)`: Creates new filesystem instance.
- `Mount(mountPoint, userStructReference)`: Mounts the filesystem to the given mount point.
