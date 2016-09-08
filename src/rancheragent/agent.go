package main

import (
        "fmt"
        "flag"
        "encoding/json"
        "net"
        "net/http"
        "time"
        "utils"
        "log"
        "strings"
)

const agent_version = "v1.0.2"

func identifyIpAddress(networkInterface string) string {
        ifaces, err := net.Interfaces()
        utils.Check(err)
        for _, i := range ifaces {
                addrs, err := i.Addrs()
                utils.Check(err)
                if i.Name == networkInterface {
                        for _, addr := range addrs {
                                switch v := addr.(type) {
                                case *net.IPNet:
                                        return v.IP.String()
                                case *net.IPAddr:
                                        return v.IP.String()
                                }
                        }
                }
        }
        panic("No IP address found for " + networkInterface)
 }

type rancherEnvironmentResponse struct {
        Data []struct {
                Id   string `json:"id"`
                Name string `json:"name"`
                Kind string `json:"kind"`
        } `json:"data"`
}

func identifyRancherEnvironment(rancherHostname, accessKey, secretKey string) string {
        for {
                rancherEnvUrl := fmt.Sprintf("http://%s/v1/accounts", rancherHostname)
                client := &http.Client{}
                req, _ := http.NewRequest("GET", rancherEnvUrl, nil)
                req.SetBasicAuth(accessKey, secretKey)
                res, err := client.Do(req)
                utils.Check(err)
                body := rancherEnvironmentResponse{}
                err = json.NewDecoder(res.Body).Decode(&body)
                utils.Check(err)

                for i := range body.Data {
                        env := body.Data[i]
                        if env.Name == "Default" && env.Kind == "project" {
                                return env.Id
                        }
                }
                log.Println("Environment not found, waiting")
                time.Sleep(5 * time.Second)
        }
}
type registrationResponse struct {
        Data []struct {
            RegistrationUrl string
        }
}

func identifyRegistrationUrl(rancherHostname, accessKey, secretKey, rancherEnvironment string) string {
        for {
                rancherTokensUrl := fmt.Sprintf("http://%s/v1/projects/%s/registrationtokens?state=active&limit=-1",
                       rancherHostname, rancherEnvironment)
                client := &http.Client{}
                req, err := http.NewRequest("GET", rancherTokensUrl, nil)
                utils.Check(err)
                req.SetBasicAuth(accessKey, secretKey)
                res, err := client.Do(req)
                utils.Check(err)
                body := registrationResponse{}
                json.NewDecoder(res.Body).Decode(&body)

                registrationUrl := body.Data[0].RegistrationUrl

                if registrationUrl != "" {
                        return registrationUrl
                }

                log.Println("Registration URL not found, waiting")
                time.Sleep(5 * time.Second)
        }
}

func installRancherAgent(ip_address, labels, rancher_url string) {
        command := fmt.Sprintf(
            `sudo docker run \
            -e CATTLE_AGENT_IP=%v \
            -e CATTLE_HOST_LABELS="%s" \
            -d --privileged --name rancher-bootstrap \
            -v /var/run/docker.sock:/var/run/docker.sock \
            -v /var/lib/rancher:/var/lib/rancher \
              rancher/agent:%s %s`,
                ip_address,
                labels,
                agent_version,
                rancher_url)
        log.Println(command)
        utils.Execute(command, nil, "")
}

func main() {
        rancherHostname := flag.String("rancher", "", "Rancher hostname")
        rancherAccessKey := flag.String("access_key", "", "Rancher Access Key")
        rancherSecretKey := flag.String("secret_key", "", "Rancher Secret Key")
        networkInterface := flag.String("interface", "", "Network Interface used to access Rancher Server")
        labelsPtr := flag.String("labels", "", "Comman separated set of NAME=VALUE labels for this host")
        flag.Parse()

        labels := strings.Replace(*labelsPtr, ",", "&", -1)

        utils.Required(*rancherHostname, "-rancher required")
        utils.Required(*rancherAccessKey, "-access_key required")
        utils.Required(*rancherSecretKey, "-secret_key required")
        utils.Required(*networkInterface, "-interface required")

        ipAddress := identifyIpAddress(*networkInterface)
        log.Println("IP address: ", ipAddress)
        rancherEnvironment := identifyRancherEnvironment(*rancherHostname, *rancherAccessKey, *rancherSecretKey)
        log.Println("Environment:", rancherEnvironment)
        registrationUrl := identifyRegistrationUrl(*rancherHostname, *rancherAccessKey, *rancherSecretKey, rancherEnvironment)
        installRancherAgent(ipAddress, labels, registrationUrl)
}