# GoSSH  -  Open Source Go Infrastucture Automation Tool

![](https://github.com/Aponiatowski/GoSSH/workflows/GoSSH/badge.svg)     [![Go Report Card](https://goreportcard.com/badge/github.com/APoniatowski/GoSSH)](https://goreportcard.com/report/github.com/APoniatowski/GoSSH)   [![codebeat badge](https://codebeat.co/badges/e53dab58-a0df-4699-a4d6-cfe67fbd9b81)](https://codebeat.co/projects/github-com-aponiatowski-gossh-master)

> CI: Forgejo Actions (`.forgejo/workflows/ci.yml`, `docker`-tagged DinD runners) runs build, gofmt, vet and `go test -race`, plus four integration suites over real SSH against throwaway systemd containers: a single-host smoke test, a 20-stage feature **matrix** (ubuntu + debian), a 3-tier **infra** scenario (nginx → api → redis), and a **breadth** suite spanning all four package managers (apt/dnf/pacman/zypper) plus a non-root sudo path.


### Current version -> **v2.0.0**

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

# Baselines (v2)

Beyond ad-hoc commands, GoSSH can enforce a declarative **baseline** — a desired
state defined in YAML — across a fleet, and verify compliance without making
changes. This is the SaltStack/Ansible/Chef/Puppet-style capability, kept fast by
the same per-server concurrent SSH engine.

```
GoSSH baseline apply  <name>    # bring servers to the desired state (./config/<name>.yml)
GoSSH baseline check  <name>    # read-only compliance check (no changes)
GoSSH baseline verify <name>    # validate the baseline file itself
```

A baseline file is `name -> server-group -> stages`, where the group name matches
a group in `pool.yml`. Stages run in order:

* **Exclude** — skip servers by OS or FQDN.
* **Prerequisites** — tools to install, files to fetch (URL/local push/mount),
  VCS, custom commands.
* **Must-Have** — `Installed`, `Configured` (config files pushed to hosts),
  `Users`, `Enabled`/`Disabled` services, firewall `Rules` (open/closed),
  `Policies`, `Mounts`. (Configured + Users run before Enable/Disable so services
  start with their config and users in place.)
* **Must-Not-Have** — packages/services/users/rules/mounts that must be absent.
* **Final** — custom commands, artifact `Collect` (logs/stats/files pulled to
  `./collections/<host>/`), service reload, reboot.

See `config/EXAMPLE-baselines.yml` for a fully-annotated example, or generate a
template with `GoSSH generate baseline template`.

**Exit codes** (`apply`/`check`): `0` compliant, `2` non-compliant (a must-have is
missing or a must-not-have is present), `1` a host was unreachable / errored — so
`baseline check` is usable as a CI/monitoring gate.

Pool entries (`./config/pool.yml`) are `group -> server -> {FQDN, Username,
Password, Key_Path, Port, OS}` (order-independent; `Key_Path`/`Port` optional).
`OS` selects the package-manager/command dialect (debian, ubuntu, centos, rhel,
fedora, opensuse, sles, arch, freebsd). A non-root `Username` makes GoSSH `sudo`
each command and feed the configured password.

# Testing

Unit tests (incl. the race-tested concurrency engine and golden parser tests):

```
go test -race ./...
```

Integration suites run in CI on Docker-in-Docker; reproduce any locally with, e.g.:

```
docker compose -f test/integration/matrix/docker-compose.yml up --build \
    --abort-on-container-exit --exit-code-from harness
docker compose -f test/integration/matrix/docker-compose.yml down -v
```

(`matrix`, `infra`, and `breadth` each have their own compose project under
`test/integration/`.)

### Outstanding issues
* The known_hosts file is causing some issues (issue open for it). Ignoring known_hosts for now.
* Baseline features deferred to VM-based testing: nfs mounts, firewall zones
  (firewalld), VCS prerequisites (currently URL-fetch, not yet real `git clone`).
