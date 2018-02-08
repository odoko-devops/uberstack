package uber

import (
	"os"
	"os/exec"
	"fmt"
	"strings"
	"errors"
	"path"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

type EnvFile map[string]string


type Uber struct {
	UberSource      string
	RancherBinary   string
	EnvFile         EnvFile
	Arguments       []string
	Services        []Service
	Action          string
	ActionArguments []string
}

type Service struct {
	UberSource string
	Stack      string
	Service    string
	Type       string
}

func (s *Service) getComposePath(composeFile string) string {
	return path.Join(s.UberSource, s.Stack, s.Service, composeFile)
}

func (u *Uber) Init(ctx *cli.Context) error {

	logrus.Debug("Init")
	filename, _ := filepath.Abs(ctx.GlobalString("config"))
	logrus.Debug("Load", filename)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var envFile EnvFile
	err = yaml.Unmarshal(yamlFile, &envFile)
	if err != nil {
		return err
	}

	logrus.Debug(envFile)
	u.EnvFile = envFile
	if uberSource, ok := envFile["uber_source"]; ok {
		u.UberSource = uberSource
		logrus.Debug("UberSource=", u.UberSource)
	} else {
		return errors.New("No uber_source provided in uber.yml")
	}
	u.RancherBinary = "rancher"
	u.ActionArguments = []string{}
	return nil
}

func (u *Uber) expandServices() error {

	logrus.Debug("expandServices")
	u.Services = []Service{}

	arguments := u.Arguments
	if arguments == nil || len(arguments) == 0 {
		arguments = []string{}
		files, err := ioutil.ReadDir(u.UberSource)
		if err != nil {
			return err
		}
		for _, file := range files {
			logrus.Debug("Checking ", file.Name(), file.IsDir())
			if file.IsDir() {
				logrus.Debug("Adding ", file.Name())
				arguments = append(arguments, file.Name())
				logrus.Debug("Arguments: ", arguments)
			}
		}
		logrus.Debug("Arguments: ", arguments)

	}
	logrus.Debug("Processing arguments: ", arguments)
	for _, argument := range arguments {
		slashCount := strings.Count(argument, "/")
		logrus.Debug("Checking out ", argument)
		logrus.Debug("Slashes", slashCount)
		if slashCount == 0 {
			services, err := u.expandStack(argument)
			if err != nil {
				return err
			}
			u.Services = append(u.Services, services...)
		} else if (slashCount == 1) {
			parts := strings.Split(argument, "/")
			stackName := parts[0]
			serviceName := parts[1]
			service := Service{
				u.UberSource,
				stackName,
				serviceName,
				"compose",
			}
			u.Services = append(u.Services, service)
		} else {
			return fmt.Errorf("Too many slashes in %s", argument)
		}
	}
	logrus.Debug(u.Services)
	return nil
}

func (u *Uber) expandStack(stackName string) ([]Service, error) {

	logrus.Debug("expandStack")
	services := []Service{}

	files, err := ioutil.ReadDir(path.Join(u.UberSource, stackName))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		serviceName := file.Name()
		service := Service{
			u.UberSource,
			stackName,
			serviceName,
			"compose",
		}
		services = append(services, service)
	}
	return services, nil
}

func (u *Uber) Execute() error {
	logrus.Debug("Execute")
	u.expandServices()
	for _, service := range u.Services {
		err := u.executeService(service)
		if err!=nil {
			return err
		}
	}
	return nil
}

func (u *Uber) executeService(service Service) error {
	logrus.Debug("executeService", service)
	args := []string{}
	args = append(args, u.Action)
	args = append(args, u.ActionArguments...)
	args = append(args, "-s", service.Stack)
	args = append(args, "-e", "uber.conf")
	args = append(args, "-f", service.getComposePath("docker-compose.yml"))
	args = append(args, "--rancher-file", service.getComposePath("rancher-compose.yml"))
	args = append(args, "-d")

	logrus.Info(u.RancherBinary)
	logrus.Info(strings.Join(args, " "))
	cmd := exec.Command(u.RancherBinary, args...)

	cmd.Env = []string{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	for k,v := range u.EnvFile {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(cmd.Env, "RANCHER_CLIENT_CONFIG=" + os.Getenv("RANCHER_CLIENT_CONFIG"))
	logrus.Debug("About to run", cmd)
	err := cmd.Run()
	return err
}
