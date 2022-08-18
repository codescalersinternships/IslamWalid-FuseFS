package fs

import (
	"fmt"
	"log"
	"reflect"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/fatih/structs"
)

type FS struct {
    UserStruct any
}

type EntryGetter interface {
	GetDirentType() fuse.DirentType
}

func NewFS(userStruct any) *FS {
    return &FS{
    	UserStruct: userStruct,
    }
}

func Mount(mountPoint string, userStruct any) error {
    conn, err := fuse.Mount(mountPoint)
    if err != nil {
        return err
    }
    defer func() {
        err := conn.Close()
        if err != nil {
            log.Println(err)
        }
        fuse.Unmount(mountPoint)
    }()

    err = fs.Serve(conn, NewFS(userStruct))
    if err != nil {
        return err
    }

    return nil
}

func (f *FS) Root() (fs.Node, error) {
    dir := NewDir()
    structMap := structs.Map(f.UserStruct)
    dir.Entries = f.createEntries(structMap)
    return dir, nil
}

func (f *FS) createEntries(structMap map[string]any) map[string]any {
    entries := map[string]any{}

    for key, val := range structMap {
        if reflect.TypeOf(val).Kind() == reflect.Map {
            dir := NewDir()
            dir.Entries = f.createEntries(val.(map[string]any))
            entries[key] = dir
        } else {
            content := []byte(fmt.Sprintln(reflect.ValueOf(val)))
            file := NewFile(key, len(content), f.UserStruct)
            entries[key] = file
        }
    }

    return entries
}
