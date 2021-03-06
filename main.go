package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitSSH "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func main() {

	url := os.Args[1]
	sshCred := os.Getenv("SSH_CRED")

	name := time.Now().UnixNano()

	err := createFirstFile(url, sshCred, name)
	checkError(err)

	fs, err := getFS(url, sshCred)
	checkError(err)

	p1 := path.Join("test-dirs", fmt.Sprintf("dir-%d", name), "fiel1.txt")
	err = fs.WriteFile(p1, "some updated data")
	checkError(err)

	p2 := path.Join("test-dirs", fmt.Sprintf("dir-%d", name), "somedir", "file2.txt")
	err = fs.WriteFile(p2, "some data 2")
	checkError(err)

	p3 := path.Join("test-dirs", fmt.Sprintf("dir-%d", name), "somedir", "anotherdir", "file3.txt")
	err = fs.WriteFile(p3, "some data 3")
	checkError(err)

	err = saveFS(sshCred, fs)
	checkError(err)

}

func createFirstFile(url, sshCred string, name int64) error {
	fs, err := getFS(url, sshCred)
	if err != nil {
		return err
	}

	p1 := path.Join("test-dirs", fmt.Sprintf("dir-%d", name), "fiel1.txt")
	err = fs.WriteFile(p1, "some data")
	if err != nil {
		return err
	}

	err = saveFS(sshCred, fs)
	if err != nil {
		return err
	}

	return nil
}

func getFS(url, sshCred string) (FileSystem, error) {
	auth, err := getAuth(sshCred)
	if err != nil {
		return nil, err
	}

	rep, err := git.CloneContext(context.TODO(), memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:  url,
		Auth: auth,
	})
	if err != nil {
		return nil, err
	}

	w, err := rep.Worktree()
	if err != nil {
		return nil, err
	}

	return &filesystem{
		Filesystem: w.Filesystem,
		rep:        rep,
	}, nil
}

func saveFS(sshCred string, fs FileSystem) error {
	fss, ok := fs.(*filesystem)
	if !ok {
		return errors.New("It's not a *filesystem")
	}

	w, err := fss.rep.Worktree()
	if err != nil {
		return err
	}

	s, err := w.Status()
	if err != nil {
		return err
	}
	fmt.Println(s)

	err = w.AddGlob(".")
	if err != nil {
		return err
	}

	s, err = w.Status()
	if err != nil {
		return err
	}
	fmt.Println(s)
	fmt.Println("--")

	_, err = w.Commit("Automatic commit", &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: "Pepito",
			When: time.Now(),
		},
	})
	if err != nil {
		return err
	}

	auth, err := getAuth(sshCred)
	if err != nil {
		return err
	}

	err = fss.rep.PushContext(context.TODO(), &git.PushOptions{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	return nil
}

func getAuth(sshCred string) (transport.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(sshCred))
	if err != nil {
		return nil, err
	}

	auth := &gitSSH.PublicKeys{
		User:   "git",
		Signer: signer,
	}

	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return auth, nil
}

type FileSystem interface {
	billy.Filesystem
	WriteFile(p, d string) error
}

type filesystem struct {
	billy.Filesystem
	rep *git.Repository
}

func (fs *filesystem) WriteFile(p, d string) error {
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
