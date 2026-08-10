package channelreaderlib

import (
	"fmt"
	"sync"

	"github.com/APoniatowski/GoSSH/yamlparser"
	"github.com/superhawk610/bar"
)

// ChannelReaderAll Function to read channel until it is closed (all servers only)
func ChannelReaderAll(channel <-chan string, wg *sync.WaitGroup) {
	successcount := 0
	barp := bar.New(yamlparser.Waittotal)
	for i := 0; i < yamlparser.Waittotal; i++ {
		for message := range channel {
			if message == "OK\n" {
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

// ChannelReaderGroups Function to read channel until it is closed (groups only)
func ChannelReaderGroups(channel <-chan string, wg *sync.WaitGroup, groupSize int) {
	successcount := 0
	barp := bar.New(groupSize)
	for message := range channel {
		if message == "OK\n" {
			successcount++
		}
		barp.Tick()
	}
	barp.Done()
	fmt.Printf("%d/%d Succeeded\n", successcount, groupSize)
}
