package fs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"bazil.org/fuse"
	"github.com/fatih/structs"
)

type File struct {
	Type       fuse.DirentType
    FileName   string
    FilePath []string
	Attributes fuse.Attr
    UserStructRef any
}


func NewFile(fileName string, filePath []string, size int, userStructRef any) *File {
	return &File{
		Type:    fuse.DT_File,
        FileName: fileName,
        FilePath: filePath,
        UserStructRef: userStructRef,
		Attributes: fuse.Attr{
			Inode: 0,
            Size: uint64(size),
			Atime: time.Now(),
			Mtime: time.Now(),
			Ctime: time.Now(),
			Mode:  0o444,
		},
	}
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = f.Attributes
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
    return f.fetchFileContent(), nil
}

func (f *File) GetDirentType() fuse.DirentType {
	return f.Type
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

func (f *File) fetchFileContent() []byte {
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
    f.Attributes.Size = uint64(len(content))

    return content
}
