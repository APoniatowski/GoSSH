package examplegenerator

const baseline string = `Example - Setup webserver:            # this will be the name of a baseline that will be displayed and logged.
  Server Group 1:                     # this is the group that it will be applied to, must be IDENTICAL to group name in pool.yml.
    Exclude:                          # this will be the exlusion list, with server name(s) or OS. Can be left empty.
      OS:                             # OS's to exclude from this baseline
        - debian
        - arch
      Servers:
        - Server 1                    # must be IDENTICAL to the name in pool.yml, not FQDN.
        - Server 3
    Prerequisites:                    # prerequisite tools, actions and checks that need to be done, before anything else.
      Tools:
        - git
        - curl
        - wget
        - nfs-utils
      Files:                          # download files/scripts/etc, via different methods
        URLs:                         # it will run wget/curl for this, make sure the tools are available in Tools or Must-Have
          - http://where/file/is
        Local:                        # local transfer via scp.
          Source: /path/to/file       # the directory you are transferring it from
          Destination: /path/to/file  # the directory you are transferring it to
        Remote:                       # copy a file from a mount specified here.
          Type: nfs                   # or other mounting method, a temporary mount will be created and removed when done
          Address: 1.2.3.4
          Username: nfsuser           # leave blank if not needed
          Password: nfspassword       # leave blank if not needed
          Source: /path/from/         # as this is only a temporary mount, it will create a temp dir to mount it in, eg. /temp/path/from/
          Destination: /path/to/      # This is the directory it should copy the contents of the mount, before unmounting
          Files:                      # files you would like to copy from the mount, specify '- All' to copy everything from the mount
            - script.sh
            - some.tar.gz
      VCS:                            # this will be git, as svn and the others are losing popularity and marketshare.
        URLs:                         # Unless enough requests were made to implement it, I will add another level here,
          - https://blablabla         # to specify what should be used.
          - https://blabla
        Execute:                      # if compilation is required, one can run the commands below,
          - command 1                 # just make sure to add it to in Tools section.
          - command 2
          - command 3
      Script: /path/to/script       # Or one can keep everything blank and run a custom bash script to do prerequisite actions.
      Commands:                     # Set of commands, in case custom configs are needed
        - command 1
        - command 2
      Clean-up: true                # clean-up of tools and downloaded urls/tools/etc. And tools in Must-Have will be ignored from clean-up
    Must-Have:                      # The servers 'must have' these configured, setting the baseline
      Installed:                    # list of tools and services that need to have been installed
        - httpd
        - firewalld
        - openssh
        - policycoreutils-python
        - git
      Enabled:                      # list of tools and services that need to have been started and enabled
        - httpd
        - firewalld
        - openssh
        - rsyslog                   # it does not need to match the installed. it will check if it is running and enabled
      Disabled:                     # Beware of conflicts in Must-Not-Have and Enabled.
        - nginx
      Configured:                   # services with config files that need to be used. these will be copied to their destinations
        httpd:
          Source:
            - /path/to/config/filesrc  # 1
            - /another/filesrc       # 2
          Destination:              ######## they need match in this order
            - /path/to/config/filedst  # 1
            - /another/filedst         # 2
        openssh:
          Source:
            - /path/to/config/filesrc
          Destination:
            - /path/to/config/filedst
      Users:                        # Create users here
        webmaster:                # the name that will be created
          Groups:
            - www
          Shell: nologin          # one can create a service account
          Home-Dir:               # no home dir will be created
          Sudoer: false
        jim:
          Groups:                 # multiple groups can be added
            - wheel
            - anothergroup
          Shell: bash             # other shells can be installed and used, if needed
          Home-Dir: /path/to/dir
          Sudoer: true
      Policies:                     # selinux/apparmor policies
        Status: Enforced            # or apparmor equivalent (enforced/complains/disabled), and that will be applied
        Import: /path/to/policy     # can be left blank, or import a policy/module stored locally, and apply it remotely
        Reboot: true                # If a reboot is required for policy changes. Especially disabling selinux/apparmor
      Rules:                        # Firewall rules
        Open:                       # open/allowed ports and which protocol to use. the order must match port to protocol
          Ports:
            - 80
            - 443
          Protocols:
            - tcp
            - tcp udp               # a space between them will add both
        Closed:                     # closed/denied ports and which protocol to use. the order must match port to protocol
          Ports:
            - 8080
          Protocols:
            - tcp
        Zones:
          - public
      Mounts:                       # needed mounts
        Mount 1:                    # name your mount in this config
          Type: nfs
          Address: 1.2.3.4
          Username: nfsuser         # leave blank if not needed
          Password: nfspassword     # leave blank if not needed
          Source: /path/from/
          Destination: /path/to/
        Mount 2:                    # name your mount in this config
          Type: smb
          Address: 2.3.4.5
          Username:                 # leave blank or remove, if not needed
          Password:                 # leave blank or remove, if not needed
          Source: /path/from/
          Destination: /path/to/
    Must-Not-Have:                  # Must-Not-Have will be a lot shorter than Must-Have, as removing is a lot faster/easier than configuring
      Installed:
        - nmap
      Enabled:
        - a-service
      Disabled:                     # Beware of conflicts in Must-Have and Enabled.
        - httpd
      Users:
        - bob
        - jane
      Rules:                        # Firewall rules, that should not be there. One can leave these blank, it is just and extra check
        Open:                       # a wrong port was specified in Must-Have, but was corrected in Must-Not-Have.
          Ports:                    # It will help with troubleshooting as well, and its an extra layer to avoid issues
            - 8080
            - 8443
          Protocols:
            - tcp udp
            - tcp
        Closed:
          Ports:
            - 80
            - 443
          Protocols:
            - tcp
            - tcp udp
        Zones:
          - public
      Mounts:                       # unmounts directories that you do not want on the servers
        - /path/to/mount1
        - /path/to/mount2
    Final:                          # any final scripts that needs to be run for scripting, or final changes
      Scripts:                      # if these are empty, then it will be ignored. or everything can be ignored
        - /path/to/script           # and these scripts/commands will be run
        - /path/to/another
      Commands:
        - command 1
        - command 2
      Collect:                      # Collect/log information before finishing, eg logs/statistics/user(s) info/file(s) info
        Logs:
          - httpd
          - sshd
        Stats:
          - cpu
          - storage
          - memory
        Files:
          - /path/to/file
          - /path/to/another
        Users: true                 # collect details of currently logged in users
      Restart:
        Services: true              # will only restart services that are in Configured section, if there are none, then it defaults to false
        Servers: false              # I believe this is obvious.`
