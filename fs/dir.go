package fs

import (
	"context"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type dir struct {
	Type       fuse.DirentType
	Attributes fuse.Attr
	Entries    map[string]any
}

// NewDir creates new empty directory
func NewDir() *dir {
	return &dir{
		Type:       fuse.DT_Dir,
		Attributes: fuse.Attr{
			Inode:  0,
			Atime:  time.Now(),
			Mtime:  time.Now(),
			Ctime:  time.Now(),
			Mode:   os.ModeDir | 0o555,
		},
		Entries: map[string]any{},
	}
}

// Attr provide the core information about the directory.
func (d *dir) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = d.Attributes
	return nil
}

// Lookup serve the kernel requests for specific file or directory.
// it fetches the file system structure for the needed file of directory.
// it returns an inode or error ENOENT if not found.
func (d *dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	node, ok := d.Entries[name]
	if ok {
		return node.(fs.Node), nil
	}
	return nil, syscall.ENOENT
}

func (d *dir) GetDirentType() fuse.DirentType {
	return d.Type
}

// ReadDirAll reads all the content of a directory.
func (d *dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var entries []fuse.Dirent

	for k, v := range d.Entries {
		var a fuse.Attr
		v.(fs.Node).Attr(ctx, &a)
		entries = append(entries, fuse.Dirent{
			Inode: a.Inode,
			Type:  v.(EntryGetter).GetDirentType(),
			Name:  k,
		})
	}
	return entries, nil
}
