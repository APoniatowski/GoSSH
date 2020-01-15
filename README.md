# GoSSH  -  GoLang SSH tool

![](https://github.com/Aponiatowski/GoSSH/workflows/GoSSH/badge.svg)     [![Go Report Card](https://goreportcard.com/badge/github.com/APoniatowski/GoSSH)](https://goreportcard.com/report/github.com/APoniatowski/GoSSH)

## Project update:
It is currently in an usable state, and can be used to execute commands in varied ways and performs well. :+1:
Logging has been implemented for SSH sessions (INFOs and ERRORs) and the output has been replaced with a progress bar. All outputs from now on, will be 
written to a log file for review. If running a command was successful, why would one want to see it and clutter the terminal. 
The logs are also rotated by date, to avoid multiple logs, if time was added to the file name.
The logging will be enchanced even further, as the project continues.

Currently working on adding the ability to run commands as sudo, and also add some security for connecting to known hosts (see issue board for clarification [Issue #7](https://github.com/APoniatowski/GoSSH/issues/7) )

I have also removed some possible features, that I was planning on implementing. [^1] 
But rather chose to implement another feature for this release. [^2]

* Windows (laptop):
##### 22 production servers (across 8 different countries):

```
> real    0m7.272s
> user    0m0.062s
> sys     0m0.046s
```

* Linux (production/staging server):
##### Tested on 24 production servers (across 8 different countries):

```
> real    0m3.276s
> user    0m0.375s
> sys     0m0.062s
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
 

[^1]: Creating an client side agent, this might possibly be added for v2.0.0 release. Not guaranteed though.
[^2]: Running a bash script with little effort. Makes things simpler, than trying to cat | gossh all, or gosh all $(cat my-script.sh), etc.
