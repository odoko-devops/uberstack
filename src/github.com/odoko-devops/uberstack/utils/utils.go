package utils

import (
	"log"
	"fmt"
	"io"
	"bufio"
	"os/exec"
	"strings"
	"github.com/kr/pty"
	"time"
	"os"
	"io/ioutil"
	"bytes"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"github.com/odoko-devops/uberstack/config"
)

/***********************************************************************
 * Execute a command, with streamed output for slow running commands
 */
func watchOutputStream(typ string, r bufio.Reader, buf *bytes.Buffer) {
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Printf("%s: %s\n", typ, line)
		if buf != nil {
			buf.Write(line)
		}
	}
}

func prepareEnvironment(env config.ExecutionEnvironment) []string {
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

func Execute(command string, env config.ExecutionEnvironment, dir string) ([]byte, error) {
	args := strings.FieldsFunc(command, func(c rune)bool {
		return c == ' ' || c == '\n' || c == '\t' || c == '\\'
	})
	cmd := exec.Command(args[0], args[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}
	fmt.Printf("Executing %s\n", command)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	outputBuffer := new(bytes.Buffer)
	stdOutReader := bufio.NewReader(stdout)
	stdErrReader := bufio.NewReader(stderr)
	go watchOutputStream("stdout", *stdOutReader, outputBuffer)
	go watchOutputStream("stderr", *stdErrReader, nil)
	cmd.Wait()
	return outputBuffer.Bytes(), nil
}


/***********************************************************************
 * Execute a command, and return the output
 */
func ExecuteAndRetrieve(command string, env config.ExecutionEnvironment, dir string) (string, error) {
	fmt.Printf("Execute (with retrieve) %s\n", command)
	args := strings.FieldsFunc(command, func(c rune)bool {
		return c == ' ' || c == '\n' || c == '\t' || c == '\\'
	})
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(output), "\n"), nil
}

func sendToPty(input string, pty *os.File) {
	time.Sleep(5 * time.Second)

	pty.Write([]byte(input+"\n"))
}

func ExecuteWithInput(command, input string, env config.ExecutionEnvironment, dir string) error {
	args := strings.FieldsFunc(command, func(c rune)bool {
		return c == ' ' || c == '\n' || c == '\t' || c == '\\'
	})
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	go sendToPty(input, f)
	go io.Copy(os.Stdout, f)
	cmd.Wait()
	return nil
}

/*
type Dependency struct {
	Version string
	Url string
	ExtractCommand string
}

var dependencies = map[string]Dependency{
	"docker-machine": {
		Version: "v0.7.0",
		Url: "https://github.com/docker/machine/releases/download/$VERSION/docker-machine-$OS-$ARCH2",
	},
	"rancher-compose": {
		Version: "v0.8.6",
		Url: "https://releases.rancher.com/compose/$VERSION/rancher-compose-$OS-$ARCH1-$VERSION.tar.gz",
		ExtractCommand: "tar -C $BINARIES --strip-components=1 -xzf /tmp/download",
	},
	"terraform": {
		Version: "0.6.16",
		Url: "https://releases.hashicorp.com/terraform/$VERSION/terraform_$VERSION_$OS_$ARCH1.zip",
		ExtractCommand: "unzip -q -d $BINARIES/ /tmp/download",
	},
}

func downloadExecutable(url, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	os.Chmod(filepath, 0755)

}

func Download(cmd string) {

	dependency := dependencies[cmd]

	dependencyPath, err := exec.LookPath(cmd)
	if err == nil {
		fmt.Printf("%s already in path\n", cmd)
	} else {
		dependencyPath = fmt.Sprintf("%s/bin/%s", GetUberState(), cmd)

		binariesPath := GetUberState() + "/bin"
		os.MkdirAll(binariesPath, 0755)
		_, err = os.Stat(dependencyPath)
		if err == nil {
			fmt.Printf("%s already downloaded\n", cmd)
		} else {
			url := strings.Replace(dependency.Url, "$VERSION", dependency.Version, -1)
			url = strings.Replace(url, "$OS", runtime.GOOS, -1)
			url = strings.Replace(url, "$ARCH1", runtime.GOARCH, -1)
			url = strings.Replace(url, "$ARCH2", "x86_64", -1)  // @todo Don't hardwire this

			fmt.Printf("Downloading %s from %s\n", cmd, url)
			if dependency.ExtractCommand == "" {
				downloadExecutable(url, dependencyPath)
			} else {
				downloadExecutable(url, "/tmp/download")
				command := strings.Replace(dependency.ExtractCommand, "$BINARIES", binariesPath, -1)
				command = strings.Replace(command, "$VERSION", dependency.Version, -1)
				fmt.Printf("Extracting...\n")
				Execute(command, nil, "")
			}
		}
	}
}
*/
func Ask(message string) string {
	var answer string
	fmt.Printf("%s: ", message)
	fmt.Scanln(&answer)
	return answer
}

func Confirm(message, expected string) bool {
	return Ask(message) == expected
}

func ReadYamlFile(filename string, obj interface{}) error {

	filepath, err := Resolve(filename, true)
	if err != nil {
		return err
	}
	if (!strings.HasSuffix(filepath, ".yml")) {
		filepath = filepath + ".yml"
	}
	log.Printf("Reading %s\n", filepath)
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, obj)
	if err != nil {
		return err
	}
	return nil
}

func Resolve(filename string, isYaml bool) (string, error) {
	var p string
	if len(filename)>0 && filename[0]=='/' {
		p = filename
	} else {
		uberHome := os.Getenv("UBER_HOME")
		p = filepath.Join(uberHome, filename)
	}
	if isYaml && !strings.HasSuffix(filename, ".yml") {
		p = p + ".yml"
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return "", fmt.Errorf("Cannot locate configuration %s at %s\n", filename, p)
  	}
	return p, nil
}
