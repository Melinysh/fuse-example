package main

import (
	"flag"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type Node struct {
	inode uint64
	name  string
}

var inode uint64
var Usage = func() {
	log.Printf("Usage of %s:\n", os.Args[0])
	log.Printf("  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func NewInode() uint64 {
	inode += 1
	return inode
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		Usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	if p := c.Protocol(); !p.HasInvalidate() {
		log.Panicln("kernel FUSE support is too old to have invalidations: version %v", p)
	}
	srv := fs.New(c, nil)
	filesys := &FS{
		&Dir{Node: Node{name: "head", inode: NewInode()}, files: &[]*File{
			&File{Node: Node{name: "hello", inode: NewInode()}, data: "hello world!"},
			&File{Node: Node{name: "aybbg", inode: NewInode()}, data: "send notes"},
		}, directories: &[]*Dir{
			&Dir{Node: Node{name: "left", inode: NewInode()}, files: &[]*File{
				&File{Node: Node{name: "yo", inode: NewInode()}, data: "ayylmaooo"},
			},
			},
			&Dir{Node: Node{name: "right", inode: NewInode()}, files: &[]*File{
				&File{Node: Node{name: "hey", inode: NewInode()}, data: "heeey, thats pretty good"},
			},
			},
		},
		}}
	log.Println("About to serve fs")
	if err := srv.Serve(filesys); err != nil {
		log.Panicln(err)
	}
	// Check if the mount process has an error to report.
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Panicln(err)
	}
}
