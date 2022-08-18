package fs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"bazil.org/fuse"
	"github.com/fatih/structs"
)

func NewFile(fileName string, size int, userStruct any) *File {
	return &File{
		Type:    fuse.DT_File,
        FileName: fileName,
        UserStruct: userStruct,
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

type File struct {
	Type       fuse.DirentType
    FileName   string
	Attributes fuse.Attr
    UserStruct any
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	*a = f.Attributes
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
    return f.fetchFile(), nil
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

func (f *File) fetchFile() []byte {
    structMap := structs.Map(f.UserStruct)
    var result []byte
    var traverse func(map[string]any)

    traverse = func(m map[string]any) {
        for key, val := range m {
            if reflect.TypeOf(val).Kind() == reflect.Map {
                traverse(val.(map[string]any))
            } else {
                if key == f.FileName {
                    result = []byte(fmt.Sprintln(reflect.ValueOf(val)))
                    f.Attributes.Size = uint64(len(result))
                }
            }
        }
    }

    traverse(structMap)

    return result
}
