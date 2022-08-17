package main

import (
	"fmt"
	"reflect"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/fatih/structs"
)

type FS struct {
    userStruct any
}

type EntryGetter interface {
	GetDirentType() fuse.DirentType
}

func NewFS(userStruct any) *FS {
    return &FS{
    	userStruct: userStruct,
    }
}

func (f *FS) Root() (fs.Node, error) {
    dir := NewDir()
    structMap := structs.Map(f.userStruct)
    dir.Entries = createEntries(structMap)
    return dir, nil
}

func createEntries(structMap any) map[string]any {
    entries := map[string]any{}

    for key, val := range structMap.(map[string]any) {
        if reflect.TypeOf(val).Kind() == reflect.Map {
            dir := NewDir()
            dir.Entries = createEntries(val)
            entries[key] = dir
        } else {
            file := NewFile([]byte(fmt.Sprint(reflect.ValueOf(val))))
            entries[key] = file
        }
    }

    return entries
}
