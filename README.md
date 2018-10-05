# How to run it

With the `SSH_CRED` beeing your creds. And that's the output error

```
$> SSH_CRED=$(cat ~/.ssh/go-git-test) go run main.go git@github.com:xescugc/go-git-test.git`
?? test-dirs/dir-1538737240573061524/fiel1.txt

A  test-dirs/dir-1538737240573061524/fiel1.txt

--
panic: malformed unpack status: 0024unpack index-pack abnormal exit

goroutine 1 [running]:
main.checkError(0x8a2f40, 0xc00036a2a0)
        /home/xescugc/go/src/github.com/xescugc/go-git-test/main.go:194 +0x86
main.main()
        /home/xescugc/go/src/github.com/xescugc/go-git-test/main.go:30 +0x124
exit status 2 

```

# Solution!

At the end was only a mater of updating to `4.7` form `4.6`
