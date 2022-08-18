## Gather

Execute multiple shell commands and gather the output in one terminal window.

![](https://drone.errorcode.io/api/badges/natewehr/gather/status.svg)

## Install
### MacOS

```
$ brew tap nwehr/tap
$ brew install gather
```
### Linux

```
$ git clone https://github.com/nwehr/gather
$ cd gather
$ make install
```
## Usage

```
$ gather --cmd "start_1.sh" --cmd "start_2.sh"

=======> start_1.sh <=======

Adding K to Every Word...
Coupling Decouplers...
Amending Laws of Physics...
Ready...

=======> start_2.sh <=======

Pressing Red Button...
Reinventing Wheel...
Combobulating Discombobulator...
Ready...

```

## Options

```  
--retries <retries>    Optional: Number of times to retry failed cmd
--retry-delay <delay>  Optional: Wait time in ms before each retry
--wait <delay>         Optional: Wait time in ms before cmd is started
--dir | -d <dir>       Optional: Working dir of cmd
--cmd | -c <cmd>
```

## YouTube

I documented the process of building `gather` on youtube. 

* [#1 gather dev vlog!](https://www.youtube.com/watch?v=s8CkL0WU1s0)
* [#2 gather dev vlog!](https://www.youtube.com/watch?v=2FIvfAAPDOg)