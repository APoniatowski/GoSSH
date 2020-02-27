package examplegenerator

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"
)

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

// GeneratePool Generates a pool.yml file, if it does not exist. It will check the OS and determine
// where to write the file
func GeneratePool() error {
	path := "./config/"
	_, errF := os.Stat(path + "pool.yml")
	if errF == nil {
		return fmt.Errorf("File exists")
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil || os.IsExist(err) {
		okFile, err := os.OpenFile(path+"pool.yml", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer okFile.Close()
		logger := log.New(okFile, "", 0)
		logger.Print(pool)
	} else {
		return fmt.Errorf("Error creating path and pool.yml")
	}
	return nil
}

// PrintPoolExample Prints a pool.yml example on the terminal/prompt
func PrintPoolExample() {
	fmt.Println()
	fmt.Println(pool)
	fmt.Println()
	fmt.Println("This is an example, of how the pool.yml should be structured.")
	fmt.Println("For now, make sure this file is in the same directory as your")
	fmt.Println("binary (or executable), in:")
	fmt.Println(color.Cyan.Sprint(" Linux:"), "./config/pool.yml")
	fmt.Println(color.Cyan.Sprint(" Windows:"), " .\\config\\pool.yml")
	fmt.Println(color.Yellow.Sprint("This is subject to change in the near future"))
}
