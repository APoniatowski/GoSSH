package channelreaderlib

import (
	"fmt"
	"sync"

	"github.com/APoniatowski/GoSSH/yamlparser"
	"github.com/superhawk610/bar"
)

// channelReaderAll Function to read channel until it is closed (all servers only)
func ChannelReaderAll(channel <-chan string, wg *sync.WaitGroup) {
	successcount := 0
	barp := bar.New(yamlparser.Waittotal)
	for i := 0; i < yamlparser.Waittotal; i++ {
		for message := range channel {
			if message == "Ok\n" {
				barp.Tick()
				successcount++
			} else {
				barp.Tick()
			}
		}
	}
	defer fmt.Printf("%d/%d Succeeded\n", successcount, yamlparser.Waittotal)
	defer barp.Done()
}

// channelReaderGroups Function to read channel until it is closed (groups only)
func ChannelReaderGroups(channel <-chan string, wg *sync.WaitGroup) {
	loopcountval := len(yamlparser.ServersPerGroup) - 1
	var totalsuccesscount int
	for i := 0; i < loopcountval; i++ {
		successcount := 0
		barp := bar.New(yamlparser.ServersPerGroup[i])
		for im := 0; im < yamlparser.ServersPerGroup[i]; im++ {
			for message := range channel {
				if message == "Ok\n" {
					barp.Tick()
					successcount++
					totalsuccesscount++
				} else {
					barp.Tick()
				}
			}
		}
		barp.Done()
		fmt.Printf("%d/%d Succeeded\n", successcount, yamlparser.ServersPerGroup[i])
	}
}
