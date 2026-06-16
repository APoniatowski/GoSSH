package yamlparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
)

// ReadBaseline Baseline reader
func ReadBaseline(bl string) string {
	cyan := color.Cyan.Render
	yellow := color.Yellow.Render
	yamlLocation, _ := filepath.Abs(bl) // remove -testing to revert changes, using custom server list
	bufRead, err := os.Open(yamlLocation)
	if err != nil {
		fmt.Println(yellow("Please check your baseline file, or generate one with:"))
		fmt.Println(cyan("gossh generate pool template"))
		fmt.Println(yellow("or display a example with:"))
		fmt.Println(cyan("gossh generate pool example"))
		os.Exit(-1)
	}
	defer bufRead.Close()

	scanner := bufio.NewScanner(bufRead)
	var baselineYaml []string

	for scanner.Scan() {
		baselineYaml = append(baselineYaml, scanner.Text())
	}
	parse := strings.Join(baselineYaml, "\n")

	return parse
}

// BaselineParse Baseline parser
func BaselineParse(baselinePath string) yaml.MapSlice {
	cyan := color.Cyan.Render
	green := color.Green.Render
	red := color.Red.Render
	yellow := color.Yellow.Render
	fmt.Println(yellow("Parsing data..."))
	data := ReadBaseline(baselinePath)
	var baseline yaml.MapSlice
	err := yaml.Unmarshal([]byte(data), &baseline)
	if err != nil {
		fmt.Println(yellow("Baseline could not be parsed, "), red("encountered some errors..."))
		fmt.Println(yellow("Please check your baseline.yml file, or generate one with:"))
		fmt.Println(cyan("gossh generate baseline template"))
		fmt.Println(yellow("or display a example with:"))
		fmt.Println(cyan("gossh generate baseline example"))
	} else {
		fmt.Println(yellow("Baseline read, "), green("no errors encountered..."))
		fmt.Println("=====================================================")
	}
	return baseline
}
