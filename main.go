package main

import (
    "bazil.org/fuse"
    "bazil.org/fuse/fs"
)

type Struct struct {
    String string
    Int int
    Bool bool
    Sub SubStruct
}

type SubStruct struct {
    Float float64
}

func main() {
    c, _ := fuse.Mount("./mnt", fuse.FSName("fusefs"), fuse.Subtype("tmpfs"))
    defer c.Close()
    s := &Struct{
        String: "name",
        Int:    88,
        Bool:   true,
        Sub:    SubStruct{
            Float: 3.14,
        },
    }
    fs.Serve(c, NewFS(s))
}
