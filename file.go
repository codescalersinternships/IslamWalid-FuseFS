package main

import (
	"context"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var _ = (fs.Node)((*File)(nil))
var _ = (fs.HandleReadAller)((*File)(nil))
var _ = (fs.NodeSetattrer)((*File)(nil))
var _ = (EntryGetter)((*File)(nil))

type File struct {
    Type       fuse.DirentType
    Content    []byte
    Attributes fuse.Attr
}

func NewFile(content []byte) *File {
    return &File{
        Type:    fuse.DT_File,
        Content: content,
        Attributes: fuse.Attr{
            Inode: 0,
            Atime: time.Now(),
            Mtime: time.Now(),
            Ctime: time.Now(),
            Mode:  0o444,
        },
    }
}

func (f *File) GetDirentType() fuse.DirentType {
	return f.Type
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = f.Attributes
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return f.Content, nil
}

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	if req.Valid.Atime() {
		f.Attributes.Atime = req.Atime
	}

	if req.Valid.Mtime() {
		f.Attributes.Mtime = req.Mtime
	}

	if req.Valid.Size() {
		f.Attributes.Size = req.Size
	}

	return nil
}
