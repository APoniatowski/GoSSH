# GoSSH  -  Open Source Go Infrastucture Automation Tool

![](https://github.com/Aponiatowski/GoSSH/workflows/GoSSH/badge.svg)     [![Go Report Card](https://goreportcard.com/badge/github.com/APoniatowski/GoSSH)](https://goreportcard.com/report/github.com/APoniatowski/GoSSH)   [![codebeat badge](https://codebeat.co/badges/e53dab58-a0df-4699-a4d6-cfe67fbd9b81)](https://codebeat.co/projects/github-com-aponiatowski-gossh-master)

## Project update:
It is currently in an usable state, and can be used to execute commands in varied ways and performs well. :+1:

Logging has been implemented for SSH sessions (INFOs and ERRORs) and the output has been replaced with a progress bar. All outputs from now on, will be 
written to a log file for review. If running a command was successful, why would one want to see it and clutter the terminal. 
The logs are also rotated by date, to avoid multiple logs, if time was added to the file name.
The logging will be enhanced even further, as the project continues.

Sudo commands are possible now. Just make sure you add the username's password to the password field, and it will be used when a password prompt should appear.

The known_hosts file is causing some issues (issue open for it).


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
GoSSH [ option ] [ subcommand ] [ command ]

Options:
* sequential, s  --Run the command sequentially on all servers in your config file
* groups, g      --Run the command on all servers per group concurrently in your config file
* all, a         --Run the command on all servers concurrently in your config file

Subcommand:
* run           --Run a bash script on your selected option (sequential/groups/all)

## Please feel free to test/use this and leave issues and comments in the issues tab.
## I will be actively working on this for the foreseeable future
