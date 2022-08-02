## Gather

`gather` will execute multiple shell commands and gather their output to display in one terminal window. It will prepend the output of each command with a number identifying which command it came from. 

```
$ gather --cmd 'ls -l' --cmd 'cat go.mod'
======== /bin/cat ========
module github.com/nwehr/gather

go 1.18
exited with code 0

======== /bin/ls ========
total 24
-rw-r--r--  1 natewehr  staff    40 Jul 30 14:01 go.mod
-rw-r--r--  1 natewehr  staff  1932 Aug  2 11:39 main.go
-rw-r--r--  1 natewehr  staff   751 Jul 31 18:34 readme.md
exited with code 0
```

## YouTube

I documented the process of building `gather` on youtube. 

* [#1 Creating your own dev tool in Go!](https://www.youtube.com/watch?v=s8CkL0WU1s0)
* [#2 Creating your own dev tool in Go!](https://www.youtube.com/watch?v=2FIvfAAPDOg)