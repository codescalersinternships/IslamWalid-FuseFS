package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var _ = (fs.Node)((*Dir)(nil))
var _ = (fs.HandleReadDirAller)((*Dir)(nil))
var _ = (EntryGetter)((*Dir)(nil))

type Dir struct {
	Type       fuse.DirentType
	Attributes fuse.Attr
	Entries    map[string]interface{}
}

func NewDir() *Dir {
	return &Dir{
		Type: fuse.DT_Dir,
		Attributes: fuse.Attr{
			Inode: 0,
			Atime: time.Now(),
			Mtime: time.Now(),
			Ctime: time.Now(),
			Mode:  os.ModeDir | 0o555,
		},
		Entries: map[string]interface{}{},
	}
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = d.Attributes
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	node, ok := d.Entries[name]
	if ok {
		return node.(fs.Node), nil
	}
	return nil, syscall.ENOENT
}

func (d *Dir) GetDirentType() fuse.DirentType {
	return d.Type
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
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
