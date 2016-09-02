package utils

import (
	"log"
	"fmt"
	"io"
	"bufio"
	"os/exec"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)


type Environment map[string]string

/***********************************************************************
 * Check and quit on errors
 */
func Check(err error) {
	if err != nil {
		log.Fatal(err);
		panic(err)
	}
}

/***********************************************************************
 * Assert that a variable exists
 */
func Required(value string, message string) {
	if value == "" {
		panic(message)
	}
}

/***********************************************************************
 * Execute a command, with streamed output for slow running commands
 */
func watchOutputStream(typ string, r bufio.Reader) {
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Printf("%s: %s\n", typ, line)
	}
}

func prepareEnvironment(env Environment) []string {
	if env != nil {
		preparedEnv := make([]string, len(env))
		i := 0
		for k,v := range env {
			preparedEnv[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}
		return preparedEnv
	} else {
		return nil
	}
}

func Execute(command string, env Environment, dir string) {
	cmd := exec.Command("bash", "-c", command)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Start()
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)
	go watchOutputStream("stdout", *stdoutReader)
	go watchOutputStream("stderr", *stderrReader)
	cmd.Wait()
}

/***********************************************************************
 * Execute a command, and return the output
 */
func ExecuteAndRetrieve(command string, env Environment, dir string) string {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.Output()
	Check(err)
	return string(output)
}


/***********************************************************************
 * Write commands to a script file for manual execution
 */
func WriteScript(path, script string) {
	err := ioutil.WriteFile(path, []byte(script), 0755)
	Check(err)
}


/***********************************************************************
 * Ask the user to take a specific action
 */
func Ask(cmd string) {
	fmt.Printf(`
  Some commands cannot be executed within a container. They have been added to
  a script, which you must now execute within your local host.

  Please execute the following command:

  %s\n`, cmd)
}

func ReadYaml(path string, data *interface{}) {
	bytes, err := ioutil.ReadFile(path)
	Check(err)
	err = yaml.Unmarshal(bytes, &data)
	Check(err)
}