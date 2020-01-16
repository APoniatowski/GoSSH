# GoSSH  -  GoLang SSH tool

![](https://github.com/Aponiatowski/GoSSH/workflows/GoSSH/badge.svg)     [![Go Report Card](https://goreportcard.com/badge/github.com/APoniatowski/GoSSH)](https://goreportcard.com/report/github.com/APoniatowski/GoSSH)

**WIP**

## Project update:
It is currently in an usable state, and can be used to execute commands in varied ways and performs well. :+1:

* Windows (laptop):
##### 22 production servers (across 8 different countries):

```
real    0m5.366s
user    0m0.062s
sys     0m0.061s
```

* Linux (production/staging server):
##### Tested on 24 production servers (across 8 different countries):

```
real    0m3.410s
user    0m0.443s
sys     0m0.051s
```

Command run:

```> GoSSH.exe all hostname```

and

```> GoSSH all hostname```



# Current usage for GoSSH:
GoSSH [ option ] [ command ]

Options:
* seq           - Run the command sequentially on all servers in your config file
* groups        - Run the command on all servers per group concurrently in your config file
* all           - Run the command on all servers concurrently in your config file

## Please feel free to test/use this and leave issues and comments in the issues tab.
## I will be actively working on this for the foreseeable future
 
