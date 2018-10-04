package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"gopkg.in/src-d/go-billy.v3/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {

	url := os.Args[1]

	rep, stderr = git.CloneContext(context.TODO(), memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})

	w, err := rep.Worktree()
	checkError(err)

	fs := w.Filesystem

	name := time.Now().UnixNano()

	p1 := path.Join(fmt.Sprintf("dir-%d", name), name, "fiel1.txt")
	err := writeFile(p1, "some data")
	checkError(err)

	p2 := path.Join(fmt.Sprintf("dir-%d", name), name, "anotherdir", "file2.txt")
	err := writeFile(p1, "some data 2")

	s, err := w.Status()
	checkError(err)
	fmt.Println(s)

	err = w.AddGlob(".")
	s, err := w.Status()
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

func writeFile(p, d string) error {
	f, err := fs.Create(p)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, d)
	if err != nil {
		return err
	}
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
