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
	"net/http"
	"runtime"
	"io/ioutil"
	"text/template"
	"bytes"
)


type Environment map[string]string

/***********************************************************************
 * Check and quit on errors
 */
func Check(err error) {
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
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
		uberState := GetUberStateLenient()
		if uberState != "" {
			env["PATH"] = "/bin:/usr/bin:/usr/local/bin/:" + uberState + "/bin"
		}
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

func splitFunc(c rune) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\\'
}

func Execute(command string, env Environment, dir string) {
	args := strings.FieldsFunc(command, splitFunc)
	cmd := exec.Command(args[0], args[1:]...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}
	fmt.Printf("Executing %s\n", command)
	cmd.Start()
	stdOutReader := bufio.NewReader(stdout)
	stdErrReader := bufio.NewReader(stderr)
	go watchOutputStream("stdout", *stdOutReader)
	go watchOutputStream("stderr", *stdErrReader)
	cmd.Wait()
}

/***********************************************************************
 * Execute a command, and return the output
 */
func ExecuteAndRetrieve(command string, env Environment, dir string) string {
	fmt.Printf("Execute (with retrieve) %s\n", command)
	args := strings.FieldsFunc(command, splitFunc)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.Output()
	Check(err)
	return strings.TrimRight(string(output), "\n")
}

/***********************************************************************
 * Execute a command on a remote Docker host
 */
func ExecuteRemote(host, cmd string, env Environment, dir string) {
	command := fmt.Sprintf(`docker-machine -s %s/machine ssh %s %s`, GetUberState(), host, cmd)
	Execute(command, env, dir)
}

func sendToPty(input string, pty *os.File) {
	time.Sleep(5 * time.Second)

	pty.Write([]byte(input+"\n"))
}

func ExecuteWithInput(command, input string, env Environment, dir string) {
	args := strings.FieldsFunc(command, splitFunc)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Env = prepareEnvironment(env)

	if dir != "" {
		cmd.Dir = dir
	}

	f, err := pty.Start(cmd)
	Check(err)

	go sendToPty(input, f)
	go io.Copy(os.Stdout, f)
	cmd.Wait()
}

func GetUberState() string {
	uberState := os.Getenv("UBER_STATE")

	if uberState == "" {

		uberHome := os.Getenv("UBER_HOME")
		if uberHome == "" {
			println("Please set either UBER_HOME or UBER_STATE")
			os.Exit(1)
		}
		uberState = uberHome + "/state"
	}
	return uberState
}

func GetUberStateLenient() string {
	uberState := os.Getenv("UBER_STATE")

	if uberState == "" {

		uberHome := os.Getenv("UBER_HOME")
		if uberHome == "" {
			return ""
		}
		uberState = uberHome + "/state"
	}
	return uberState
}

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

func downloadExecutable(url, filepath string) {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0755)
	Check(err)
	defer f.Close()

	resp, err := http.Get(url)
	Check(err)
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	Check(err)
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


func TerraformExport(config, providerName, name string, params map[string]string) {
	dir := fmt.Sprintf("%s/terraform/%s", GetUberState(), providerName)
	err := os.MkdirAll(dir, 0755)
	Check(err)

	path := fmt.Sprintf("%s/%s", dir, name)


	configTemplate, err := template.New("terraform").Parse(config)
	Check(err)

	buf := bytes.Buffer{}
	configTemplate.Execute(&buf, params)

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	Check(err)
}

func TerraformApply(providerName string, resources []string, env Environment) {
	path := fmt.Sprintf("%s/terraform/%s", GetUberState(), providerName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Cannot apply terraform config - cannot find %s\n", path)
		os.Exit(1)
  	}

	resourceTargets := make([]string, len(resources))
	for i, resource := range resources {
		resourceTargets[i] = "-target=" + resource
	}
	command := "terraform apply -refresh=true " + strings.Join(resourceTargets, " ")
	Execute(command, env, path)
}

func TerraformOutput(providerName string, output string) string {
	path := fmt.Sprintf("%s/terraform/%s", GetUberState(), providerName)
	command := fmt.Sprintf("terraform output %s", output)
	return ExecuteAndRetrieve(command, nil, path)
}

func TerraformDestroy(name string, env Environment) {
	path := fmt.Sprintf("%s/terraform/%s", GetUberState(), name)
	command := fmt.Sprintf("terraform destroy -state=%s/terraform.tfstate -force", path)
	Execute(command, env, path)
}

func TerraformRemoveState(providerName string) {
	uberState := GetUberState()
	path1 := fmt.Sprintf("%s/terraform/%s/terraform.tfstate", uberState, providerName)
	path2 := fmt.Sprintf("%s/terraform/%s/terraform.tfstate.backup", uberState, providerName)
	os.Remove(path1)
	os.Remove(path2)
}