package examplegenerator

const pool string = `ServerGroup1:     #group name, spacing does not matter
  Server11:       #server name, spacing does not matter
    FQDN: hostname11.whatever.com #
    Username: user11              ##
    Password: password11          ###     FQDN, is needed. Username defaults to root,
    Key_Path: /path/to/key        ##      password or key needed, ports default to 22
    Port: 22                      #       and OS will include all package managers per OS
    OS: debian
  Server12:
    FQDN: hostname12.whatever.com
    Username: user12
    Password: password12
    Key_Path: /path/to/key
    Port: 223
    OS: centos    # or rhel
ServerGroup2:
  Server21:
    FQDN: hostname21.whatever.com
    Username: user21
    Password: password21
    Key_Path: /path/to/key
    Port: 2233
    OS: fedora
  Server22:
    FQDN: hostname22.whatever.com
    Username: user22
    Password: password22
    Key_Path: /path/to/key
    Port:
    OS: opensuse  # or sles`
