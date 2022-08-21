# Gather

![](https://drone.errorcode.io/api/badges/natewehr/gather/status.svg)

Execute multiple shell commands and gather the output in one terminal window. 


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

# Install
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


# Options

```  
--retries <retries>    Optional: Number of times to retry failed cmd
--retry-delay <delay>  Optional: Wait time in ms before each retry
--wait <delay>         Optional: Wait time in ms before cmd is started
--dir | -d <dir>       Optional: Working dir of cmd
--cmd | -c <cmd>
```

# Donate

Bitcoin (BTC)

```
bc1qkm8gm3ggu8s4lnnc8mp0fahksp23u965hp758c
```

Ravencoin (RVN)

```
RSm7jfUjynsVptGyEDaW5yShiXbKBPsHNm
```