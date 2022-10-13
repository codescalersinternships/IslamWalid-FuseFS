package fs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"bazil.org/fuse"
	"github.com/fatih/structs"
)

type file struct {
	Type            fuse.DirentType
    FileName        string
    FilePath        []string
    Content         []byte
	Attributes      fuse.Attr
    UserStructRef   any
}

// NewFile creates new empty file.
func NewFile(fileName string, filePath []string, size int, userStructRef any) *file {
	return &file{
		Type:           fuse.DT_File,
        FileName:       fileName,
        FilePath:       filePath,
        UserStructRef:  userStructRef,
		Attributes:     fuse.Attr{
			Inode:      0,
            Size:       uint64(size),
			Atime:      time.Now(),
			Mtime:      time.Now(),
			Ctime:      time.Now(),
			Mode:       0o444,
		},
	}
}

// Attr provide the core information about the file.
func (f *file) Attr(ctx context.Context, a *fuse.Attr) error {
    f.updateFile()
	*a = f.Attributes
	return nil
}

// ReadAll returns all the content in a file.
func (f *file) ReadAll(ctx context.Context) ([]byte, error) {
    f.updateFile()
    f.Attributes.Atime = time.Now() 
    return f.Content, nil
}

func (f *file) GetDirentType() fuse.DirentType {
	return f.Type
}

func (f *file) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
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

// updateFile fetch the given file for its content and attributes.
func (f *file) updateFile() {
    var content []byte
    var traverse func(m map[string]any, idx int)

    structMap := structs.Map(f.UserStructRef)
    
    traverse = func(m map[string]any, idx int) {
        if idx == len(f.FilePath) {
            content = []byte(fmt.Sprintln(reflect.ValueOf(m[f.FileName])))
        } else {
            traverse(m[f.FilePath[idx]].(map[string]any), idx + 1)
        }
    }

    traverse(structMap, 0)

    f.Content = content
    f.Attributes.Size = uint64(len(content))
}
