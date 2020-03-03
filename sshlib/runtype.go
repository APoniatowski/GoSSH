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

/////////////////////////////////////////////////////////////////////////////////////////////////
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
			defaulter(&pp)
			go connectAndRun(command, servername.(string), &pp, output, &wg)
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
		defaulter(&pp)

		go connectAndRun(command, servername.(string), &pp, output, &wg)
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
			defaulter(&pp)
			s := spinner.New(spinner.CharSets[9], 25*time.Millisecond)
			s.Prefix = servername.(string) + ": "
			s.Start()
			output := connectAndRunSeq(command, servername.(string), &pp)
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

/////////////////////////////////////////////////////////////////////////////////////////////////
// ApplyBaselines testing
func ApplyBaselines(baselineyaml *yaml.MapSlice) {
	for _, groupItem := range *baselineyaml {
		fmt.Printf("Processing %s:\n", groupItem.Key)
		groupValue, ok := groupItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
		}
		for _, serverItem := range groupValue {
			servername := serverItem.Key
			fmt.Println(servername)
		}
	}
}

// ApplyBaselines testing
func CheckBaselines(baselineyaml *yaml.MapSlice) {
	var blstruct ParsedBaseline
	var servergroupnames []string
	// first - BL names
	for _, blItem := range *baselineyaml {
		fmt.Printf("%s:\n", blItem.Key)
		groupValues, ok := blItem.Value.(yaml.MapSlice)
		if !ok {
			panic(fmt.Sprintf("Check your baseline for issues\nAlternatively generate a template to see what is missing/wrong\n"))
		}
		// second - Server group names
		for _, groupItem := range groupValues {
			servergroupnames = append(servergroupnames, groupItem.Key.(string)) // done
			fmt.Printf("%s:\n", groupItem.Key)
			blstepsValue, ok := groupItem.Value.(yaml.MapSlice)
			if !ok {
				panic(fmt.Sprintf("Error reading Server Groups. Aborting to prevent possible damage"))
			}
			// third - BL steps or phases (Excludes, Prerequisites, Must-Haves, Must-Not-Haves,etc)
			for _, blstepItem := range blstepsValue {
				nextValues, ok := blstepItem.Value.(yaml.MapSlice)
				time.Sleep(1 * time.Second)
				if !ok {
					// If excludes/prereqs/etc are missing or empty, create empty/blank data dor structs
					// skipping those steps. An extra error will be created if too many fields are missing
				}
				fmt.Println(blstepItem.Key)
				if blstepItem.Key == nil {
					fmt.Println("blank this step")
				}
				// fourth - OS, Servers, Tools, Files, VCS, etc
				for _, thirdStep := range nextValues {
					nnnextValue, ok := thirdStep.Value.(yaml.MapSlice)
					if !ok {
						// If excludes/prereqs/etc are missing or empty, create empty/blank data dor structs
						// skipping those steps. An extra error will be created if too many fields are missing
					}
					if thirdStep.Key == "OS" {
						exclOS := make([]string, len(thirdStep.Value.([]interface{})))
						OSslice := thirdStep.Value.([]interface{})
						for i, v := range OSslice {
							exclOS[i] = v.(string)
						}
						blstruct.exclude.osExcl = exclOS
						fmt.Println("OS stored")
						time.Sleep(1 * time.Second)
					}
					if thirdStep.Key == "Servers" {
						exclServers := make([]string, len(thirdStep.Value.([]interface{})))
						// var test []interface{}
						// fmt.Println("1 ", thirdStep.Key)
						// fmt.Println("2 ", nnnextValue)
						// fmt.Println("3 ", thirdStep.Value.([]interface{}))
						serverSlice := thirdStep.Value.([]interface{})
						for i, v := range serverSlice {
							exclServers[i] = v.(string)
						}
						blstruct.exclude.serversExcl = exclServers
						fmt.Println("servers stored")
						time.Sleep(1 * time.Second)
					}
					// test2 := strings.Split(test, " ")
					//  fifth
					for _, ffforItem := range nnnextValue {
						// fmt.Println("4 ", ffforItem.Value)
						// ffforname := ffforItem.Value
						// fffornamekey := ffforItem.Key
						// fmt.Println("\t", ffforname)
						// fmt.Println(fffornamekey)
						nnnnextValue, ok := ffforItem.Value.(yaml.MapSlice)
						if !ok {
							// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
						}
						// sixth
						for _, fffforItem := range nnnnextValue {
							// fmt.Printf("%s:\n", fffforItem.Key)
							// fffforname := fffforItem.Value
							// ffffornamekey := fffforItem.Key
							// fmt.Println("\t", fffforname)
							// fmt.Println(ffffornamekey)
							nnnnnextValue, ok := fffforItem.Value.(yaml.MapSlice)
							if !ok {
								// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
							}
							for _, ffffforItem := range nnnnnextValue {
								// fmt.Printf("%s:\n", ffffforItem.Key)
								// somevalue := ffffforItem.Value
								// fmt.Println("\t", somevalue)
								_, ok := ffffforItem.Value.(yaml.MapSlice)
								if !ok {
									// panic(fmt.Sprintf("Unexpected type %T", groupItem.Value))
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Println(servergroupnames)
	fmt.Println(blstruct.exclude.osExcl)
	fmt.Println(blstruct.exclude.serversExcl)
}
