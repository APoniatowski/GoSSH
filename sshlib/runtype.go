package sshlib

import (
	"fmt"
	"sync"
	"time"

	"github.com/APoniatowski/GoSSH/channelreaderlib"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

// RunGroups This will run servers concurrently and groups sequentially
func RunGroups(configs *yaml.MapSlice, command *string) {
	for _, groupItem := range *configs {
		output := make(chan string)
		var wg sync.WaitGroup
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			wg.Add(1)
			servername := serverItem.Key
			serverValue, ok := serverItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
			}
			var pp ParsedPool
			pp.fqdn = serverValue[0].Value
			pp.username = serverValue[1].Value
			pp.password = serverValue[2].Value
			pp.keypath = serverValue[3].Value
			pp.port = serverValue[4].Value
			pp.os = serverValue[5].Value
			pp.defaulter()
			go pp.connectAndRun(command, servername.(string), output, &wg)
		}
		go func() {
			wg.Wait()
			close(output)
		}()
		channelreaderlib.ChannelReaderGroups(output, &wg)
	}

}

// RunAllServers As the function implies, this will run all servers concurrently
func RunAllServers(configs *yaml.MapSlice, command *string) {
	var allServers yaml.MapSlice
	output := make(chan string)
	var wg sync.WaitGroup
	// Concatenates the groups to create a single group
	for _, groupItem := range *configs {
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		allServers = append(allServers, groupValue...)
	}
	for _, serverItem := range allServers {
		wg.Add(1)
		servername := serverItem.Key
		serverValue, ok := serverItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
		}
		var pp ParsedPool
		pp.fqdn = serverValue[0].Value
		pp.username = serverValue[1].Value
		pp.password = serverValue[2].Value
		pp.keypath = serverValue[3].Value
		pp.port = serverValue[4].Value
		pp.os = serverValue[5].Value
		pp.defaulter()
		go pp.connectAndRun(command, servername.(string), output, &wg)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	channelreaderlib.ChannelReaderAll(output, &wg)
}

// RunSequentially Function for running everything sequentially, this will be the default behaviour
func RunSequentially(configs *yaml.MapSlice, command *string) {
	for _, groupItem := range *configs {
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			servername := serverItem.Key
			serverValue, ok := serverItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Unexpected type %T", serverItem.Value))
			}
			var pp ParsedPool
			pp.fqdn = serverValue[0].Value
			pp.username = serverValue[1].Value
			pp.password = serverValue[2].Value
			pp.keypath = serverValue[3].Value
			pp.port = serverValue[4].Value
			pp.os = serverValue[5].Value
			pp.defaulter()
			s := spinner.New(spinner.CharSets[9], 25*time.Millisecond)
			s.Prefix = servername.(string) + ": "
			s.Start()
			output := pp.connectAndRunSeq(command, servername.(string))
			if output == "OK\n" {
				s.Stop()
				fmt.Printf("%v: ", servername)
				fmt.Print(color.Green.Sprint(output))
			} else {
				s.Stop()
				fmt.Printf("%v: ", servername)
				fmt.Print(color.Red.Sprint(output))
			}
		}
	}
}
