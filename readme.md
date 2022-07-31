## Gather

`gather` will execute multiple shell commands and gather their output to display in one terminal window. It will prepend the output of each command with a number identifying which command it came from. 

```
$ gather --cmd 'ls -l' --cmd 'cat go.mod'
1: module github.com/nwehr/gather
1: 
1: go 1.18
0: total 24
0: -rw-r--r--  1 natewehr  staff    40 Jul 30 14:01 go.mod
0: -rw-r--r--  1 natewehr  staff  1272 Jul 31 12:54 main.go
0: -rw-r--r--  1 natewehr  staff   555 Jul 31 18:27 readme.md
```

## YouTube

I documented the process of building `gather` on youtube. 

* [#1 Creating your own dev tool in Go!](https://www.youtube.com/watch?v=s8CkL0WU1s0)
* [#2 Creating your own dev tool in Go!](https://www.youtube.com/watch?v=2FIvfAAPDOg)