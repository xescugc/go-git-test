package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {

	url := os.Args[1]

	rep, err := git.CloneContext(context.TODO(), memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})
	checkError(err)

	w, err := rep.Worktree()
	checkError(err)

	fs := w.Filesystem

	name := time.Now().UnixNano()

	p1 := path.Join(fmt.Sprintf("dir-%d", name), "fiel1.txt")
	err = writeFile(fs, p1, "some data")
	checkError(err)

	p2 := path.Join(fmt.Sprintf("dir-%d", name), "somedir", "file2.txt")
	err = writeFile(fs, p2, "some data 2")
	checkError(err)

	p3 := path.Join(fmt.Sprintf("dir-%d", name), "somedir", "anotherdir", "file3.txt")
	err = writeFile(fs, p3, "some data 3")
	checkError(err)

	s, err := w.Status()
	checkError(err)
	fmt.Println(s)

	err = w.AddGlob(".")
	s, err = w.Status()
	checkError(err)
	fmt.Println(s)

	_, err = w.Commit("Automatic commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: "Pepito",
			When: time.Now(),
		},
	})

	err = rep.PushContext(context.TODO(), &git.PushOptions{})
	checkError(err)

}

func writeFile(fs billy.Filesystem, p, d string) error {
	f, err := fs.Create(p)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, d)
	if err != nil {
		return err
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
