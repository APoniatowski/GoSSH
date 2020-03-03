package examplegenerator

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"
)

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

func GenerateBaseline() error {
	path := "./config/"
	_, errF := os.Stat(path + "baseline.yml")
	if errF == nil {
		return fmt.Errorf("File exists")
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err == nil || os.IsExist(err) {
		okFile, err := os.OpenFile(path+"baseline.yml", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer okFile.Close()
		logger := log.New(okFile, "", 0)
		logger.Print(baseline)
	} else {
		return fmt.Errorf("Error creating path and baseline.yml")
	}
	return nil
}

func PrintBaselineExample() {
	fmt.Println()
	fmt.Println(baseline)
	fmt.Println()
	fmt.Println("This is an example, of how the pool.yml should be structured.")
	fmt.Println("For now, make sure this file is in the same directory as your")
	fmt.Println("binary (or executable), in:")
	fmt.Println(color.Cyan.Sprint(" Linux:"), "./config/pool.yml")
	fmt.Println(color.Cyan.Sprint(" Windows:"), " .\\config\\pool.yml")
	fmt.Println(color.Yellow.Sprint("This is subject to change in the near future"))
}
