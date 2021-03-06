package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

//
var tag, version = "0.1", "0"

var debugging = os.Getenv("DEBUG") != ""

//
func main() {
	jq, err := exec.LookPath("jq")
	if err != nil {
		log.Fatal(err)
	}

	// This isn't ideal but, working from right to left, file the first parameter that isn't a file
	// and assume that's the filter. Everything from the filter to the left inclusively are the options to jq.

	files := []string{}
	args := os.Args[1:]
	a := len(args) - 1
	if a < 1 {
		log.Fatal("usage: yq query -- file...")
	}

	for ; a > 0; a-- {
		arg := args[a]
		if arg == "--" {
			break
		}
		files = append([]string{args[a]}, files...)
	}

	{
		stat, _ := os.Stdin.Stat()
		isTTY := (stat.Mode() & os.ModeCharDevice) != 0
		if len(files) == 0 || !isTTY {
			files = append(files, "")
		}
	}

	//

	cmd := exec.Command(jq, args[:a]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	//

	in := make(chan []byte, 0)

	go func() {
		defer stdin.Close()
		for s := range in {
			stdin.Write(s)
		}
	}()

	//

	for _, arg := range files {
		out := map[string]interface{}{}
		var file []byte
		if arg == "" {
			in, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Println("reading", "stdin", err)
				continue
			}
			file = in
		} else {
			in, err := ioutil.ReadFile(arg)
			if err != nil {
				log.Println("reading", arg, err)
				continue
			}
			file = in
		}
		err = yaml.Unmarshal(file, &out)
		if err != nil {
			log.Println("yaml", arg, err)
			continue
		}
		s, err := json.Marshal(fix(out))
		if err != nil {
			log.Println("json", arg, err)
		} else {
			if (debugging) {
				log.Println("jq input\n", string(s))
			}
			in <- s
		}
	}
	close(in)

	// Run the command and get the json output

	jout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(jout), err)
		os.Exit(4)
	}

	if (debugging) {
		log.Println("jq output\n", string(jout))
	}

	// Convert the jq output back to yaml

	var yout interface{}
	err = json.Unmarshal(jout, &yout)
	if err != nil {
		fmt.Println(string(jout))
		log.Println(err)
		os.Exit(3)
	}

	out, err := yaml.Marshal(yout)
	if err != nil {
		fmt.Println(string(jout))
		log.Println(err)
		os.Exit(2)
	}

	fmt.Printf("%s\n", out)
}

// Starting at the top-level map, recurse the broken map and convert all keys to strings.
func fix(out map[string]interface{}) map[string]interface{} {
	var fixs func(out []interface{}) []interface{}
	var fixm func(out map[interface{}]interface{}) map[string]interface{}

	assign := func(i string, v interface{}, fixed map[string]interface{}) {
		switch w := v.(type) {
		case []interface{}:
			fixed[i] = fixs(w)
		case map[interface{}]interface{}:
			fixed[i] = fixm(w)
		default:
			fixed[i] = w
		}
	}

	// fix map keys
	fixm = func(out map[interface{}]interface{}) map[string]interface{} {
		fixed := map[string]interface{}{}
		for k, v := range out {
			assign(fmt.Sprintf("%v", k), v, fixed)
		}
		return fixed
	}

	// fix slice elements.
	fixs = func(out []interface{}) []interface{} {
		lo := len(out)
		fixed := make([]interface{}, lo)
		for i, v := range out {
			switch w := v.(type) {
			case []interface{}:
				fixed[i] = fixs(w)
			case map[interface{}]interface{}:
				fixed[i] = fixm(w)
			default:
				fixed[i] = w
			}
		}
		return fixed
	}

	fixed := map[string]interface{}{}
	for k, v := range out {
		assign(k, v, fixed)
	}
	return fixed
}
