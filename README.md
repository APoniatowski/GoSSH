# GoSSH  -  Open Source Go Infrastucture Automation Tool

![](https://github.com/Aponiatowski/GoSSH/workflows/GoSSH/badge.svg)     [![Go Report Card](https://goreportcard.com/badge/github.com/APoniatowski/GoSSH)](https://goreportcard.com/report/github.com/APoniatowski/GoSSH)   [![codebeat badge](https://codebeat.co/badges/e53dab58-a0df-4699-a4d6-cfe67fbd9b81)](https://codebeat.co/projects/github-com-aponiatowski-gossh-master)


### Current version -> **v1.4.0**

### Goal with this project:
I've seen so many times that other tools like ansible, saltstack, etc, perform really well and give in-depth information. To any engineer (devops, IT, software)
would be extremely useful. But the common complaints I have heard (and seen it for myself), was the speed at which it does its job. I took it upon myself to learn Go
and create an useful tool to (*hopefully*) replace those others, as there are only 4 commonly used tools out there (Ansible, Saltstack, Chef and Puppet) and got tired
of the vendor lock-in, with a slow performing tool. Or one that is rediculously complex to configure.

So I went with the K.I.S.S. method, and keep the complexity in the code, not the tool. And boost the performance with a modern language.


* Windows (laptop):
##### Tested on 22 production servers (across 8 different countries):

```
 (████████████████████) 100.0% 6.2 ops/s
22/22 Succeeded

real    0m3.775s
user    0m0.061s
sys     0m0.031s
```


* Linux (production/staging server):
##### Tested on 24 production servers (across 8 different countries):

```
 (████████████████████) 100.0% 6.7 ops/s
24/24 Succeeded

real    0m3.468s
user    0m0.430s
sys     0m0.066s
```


* Linux (production/staging server):
##### Tested on 75 pre-production servers (in the same site):
```
 (████████████████████) 100.0% 172.5 ops/s
75/75 Succeeded

real    0m0.455s
user    0m0.479s
sys     0m0.219s
```

Command run:

```> GoSSH.exe all hostname```

and

```> GoSSH all hostname```


And to my surprise, this tool outperformed saltstack (probably Ansible too). I would love to get benchmarks for the other tools. Saltstack took around 3.4 seconds 
to execute the same command (`hostname`) on the same set of servers.  I wish I could test this in a bigger environment, as the one I tested it on, was the 
pre-production servers I was allowed to test it on.


##### Note:
Logs will be written to ```./logs/*``` in their individual directories (```/errors``` and ```/output```) in the same directory as where the application is used.  
Make sure the pool.yml file is in ```./config``` and saved as ```pool.yml``` 
(please use the config file in this repo as a template)


# Current usage for GoSSH:
GoSSH [ option ] [ subcommand ] [ command ]

Options:
* sequential, s  --Run the command sequentially on all servers in your pool
* groups, g      --Run the command on all servers per group concurrently in your pool
* all, a         --Run the command on all servers concurrently in your pool

Subcommand:
* run           --Run a bash script on your selected option (sequential/groups/all)
* update        --Update all packages on servers in your pool (optional os or OS flag will do a system upgrade)
* install       --Install packages on servers in your pool
* uninstall     --Uninstall packages on servers in your pool

### Outstanding issues
*The known_hosts file is causing some issues (issue open for it). Ignoring known_hosts for now.
