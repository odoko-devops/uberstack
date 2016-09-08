package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"flag"
	"gopkg.in/yaml.v2"
	"strings"
	"path/filepath"
	"utils"
	"installer/model"
)

/***********************************************************************
 * Uberstack type definitions
 */

type uberstack_type struct {
	Name         string
	Stacks       []string
	Uberstacks   []string
	Required     []string
	Environments map[string]utils.Environment
}

/***********************************************************************
 * Identify stacks within $UBER_HOME directory
 */
func get_uberstacks(uber_home string) []string {
	files, _ := ioutil.ReadDir(uber_home + "/uberstacks")
	uberstacks := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
			s := strings.Split(f.Name(), ".")
			if len(s) == 2 && s[1] == "yml" {
				uberstacks = uberstacks[0: len(uberstacks) + 1]
				uberstacks[len(uberstacks) - 1] = s[0]
			}
		}
	}
	return uberstacks
}

/***********************************************************************
 * Identify stacks within $UBER_HOME directory
 */
func get_stacks(uber_home string) []string {
	files, _ := ioutil.ReadDir(uber_home + "/stacks")
	stacks := make([]string, 0, len(files))
	for _, f := range files {
		s := strings.Split(f.Name(), ".")
		if len(s) == 1 {
			stacks = stacks[0: len(stacks) + 1]
			stacks[len(stacks) - 1] = s[0]
		}
	}
	return stacks
}

/***********************************************************************
 * Read a single uberstack from its config yaml file
 */
func get_uberstack(uber_home string, name string) uberstack_type {
	bytes, err := ioutil.ReadFile(uber_home + "/uberstacks/" + name + ".yml")
	utils.Check(err)
	uberstack := uberstack_type{}
	err = yaml.Unmarshal(bytes, &uberstack)
	utils.Check(err)
	return uberstack
}

/***********************************************************************
 * Execute ls command
 */
func ls(uber_home string) {
	fmt.Println("\nStacks:")
	fmt.Println("-------")
	stacks := get_stacks(uber_home)
	for stack_id := range stacks {
		fmt.Println(stacks[stack_id])
	}
	fmt.Println("\nUberstacks:")
	fmt.Println("-----------")
	uberstacks := get_uberstacks(uber_home)
	for uber_id := range uberstacks {
		fmt.Println(uberstacks[uber_id])
	}
}

/***********************************************************************
 * Build a suitable environment for execution
 */
func get_parameters_for(uberstack uberstack_type, env string, state_file string) utils.Environment {
	params := get_parameters_from_environment_and_uberstack(uberstack, env)
	add_parameters_from_state(env, state_file, &params)
	check_required(uberstack, params)
	return params
}

func get_parameters_from_environment_and_uberstack(uberstack uberstack_type, env string) utils.Environment {
	environ := os.Environ()
	params := utils.Environment{}

	for _, v := range environ {
		s := strings.SplitN(v, "=", 2)
		name := s[0]
		value := s[1]
		params[name] = value
	}

	for k, v := range uberstack.Environments[env] {
		params[k] = v
	}

	return params
}

func add_parameters_from_state(env string, state_file string, params *utils.Environment) {
	state := model.LoadState("/state/state.yml")

	if env == "local" {
		(*params)["RANCHER_ACCESS_KEY"] = state.Provider["virtualbox"].AccessKey
		(*params)["RANCHER_SECRET_KEY"] = state.Provider["virtualbox"].SecretKey
	} else {
		(*params)["RANCHER_ACCESS_KEY"] = state.Provider["amazonec2"].AccessKey
		(*params)["RANCHER_SECRET_KEY"] = state.Provider["amazonec2"].SecretKey
	}
}

/***********************************************************************
 * Check for required variables
 */
func check_required(uberstack uberstack_type, params utils.Environment) {

	for i := range uberstack.Required {
		required := uberstack.Required[i]
		_, ok := params[required]
		if !ok {
			log.Fatal("Required parameter: ", required)
			os.Exit(1)
		}
	}
}

/***********************************************************************
 * Process any referenced Uberstacks
 */
func process_uberstack(uber_home string, uberstack uberstack_type, env string, cmd string, exclude_stack string) {

	fmt.Println("process_uberstack", uberstack.Name)
	for i := 0; i < len(uberstack.Uberstacks); i++ {
		name := uberstack.Uberstacks[i]
		inner_uberstack := get_uberstack(uber_home, name)
		process_uberstack(uber_home, inner_uberstack, env, cmd, exclude_stack)
	}

	for i := range uberstack.Stacks {
		name := uberstack.Stacks[i]
		if name == exclude_stack {
			continue
		}
		project := name
		stack := name

		s := strings.SplitN(name, ":", 2)
		if len(s) == 2 {
			project = s[0]
			stack = s[1]
		}
		command := fmt.Sprintf(`rancher-compose --file %s/stacks/%s/docker-compose.yml \
                        --rancher-file %s/stacks/%s/rancher-compose.yml \
                        --project-name %s \
                        %s`,
			uber_home, stack, uber_home, stack, project, cmd)
		env := get_parameters_for(uberstack, env, "/state/state.yml")
		utils.Execute(command, env, "")
	}
}

func main() {

	uber_home := os.Getenv("UBER_HOME")
	if uber_home == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		utils.Check(err)
		uber_home = dir
	}

	rancher_url := os.Getenv("RANCHER_URL")
	fmt.Printf("Using UBER_HOME=%v\n", uber_home)
	fmt.Printf("Using RANCHER_URL=%v\n", rancher_url)

	excludePtr := flag.String("exclude", "", "Exclude stack")
	flag.Parse()

	/*
	if len(flag.Args()) < 3 {
	  flag.Usage()
	  os.Exit(1)
	}
	*/
	args := flag.Args()
	arg_count := len(args)
	if arg_count == 0 {
		flag.Usage()
		os.Exit(1)
	}
	action := args[0]
	uberstack_name := ""
	environment := ""
	if len(args) > 1 {
		uberstack_name = args[1]
	}
	if len(args) > 2 {
		environment = args[2]
	} else {
		environment = "local"
	}

	cmd := ""
	desc := ""
	switch action {
	case "up":
		cmd = "up -d"
		desc = "Installing"
	case "upgrade":
		cmd = "up --upgrade --pull -d " + strings.Join(flag.Args()[3:], " ")
		desc = "Upgrading"
	case "confirm-upgrade":
		cmd = "up --upgrade --confirm-upgrade"
		desc = "Confirming"
	case "rollback":
		cmd = "up --upgrade --rollback"
		desc = "Rolling back"
	case "rm":
		if environment != "local" {
			var answer string
			fmt.Print("Retype uberstack name to confirm deletion: ")
			fmt.Scanln(&answer)
			if answer != uberstack_name {
				fmt.Println("Confirmation failed, quitting")
				os.Exit(1)
			}
		}
		cmd = "rm --force"
		desc = "Removing"
	case "ls":
		ls(uber_home)
		os.Exit(0)
	default:
		fmt.Printf("Unknown action: %s", action)
		os.Exit(1)
	}
	fmt.Println("Exclude:", *excludePtr)
	fmt.Println("cmd:", cmd)
	fmt.Println("desc:", desc)
	fmt.Println("uberstack:", uberstack_name)
	uberstack := get_uberstack(uber_home, uberstack_name)
	process_uberstack(uber_home, uberstack, environment, cmd, *excludePtr)
}
