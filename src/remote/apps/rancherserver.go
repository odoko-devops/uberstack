package apps

import (
	"installer/model"
	"utils"
	"net/http"
	"time"
	"bytes"
	"fmt"
	"encoding/json"
	"log"
	"io"
)

func Rancher_InstallServer(config model.Config, state *model.State, hostConfig model.HostConfig, app model.AppConfig) {
	appHost := hostConfig.Name

	rancherHost := app.Config["host"]

	dockerHost := app.Config["docker-host"]
	realmName := app.Config["auth-realm"]
	realm := model.GetAuthRealm(config, realmName)
	email := realm.Users[0].Email
	username := realm.Users[0].Username
	password := realm.Users[0].Password

	//rancher_DockerRm(appHost) // Only use this when testing on the same host
	rancher_DockerRun(appHost)

	envId := rancher_GetEnvironmentId(rancherHost)
	accessKey, secretKey := rancher_GetApiKeys(rancherHost, envId)

	rancher_SetApiHost(rancherHost)
	registryId := rancher_RegisterRegistry(rancherHost, envId, dockerHost)
	rancher_RegistryCredentials(rancherHost, envId, registryId, username, password, email)

	rancher_EnableAuth(rancherHost, username, password)

	providerState := state.Provider[hostConfig.Provider]
	providerState.RancherUrl = rancherHost
	providerState.AccessKey = accessKey
	providerState.SecretKey = secretKey
	state.Provider[hostConfig.Provider] = providerState
	log.Printf("access: %s secret: %s", accessKey, secretKey)
}

func rancher_DockerRm(appHost string) {
	command := "docker rm -f \\$\\(docker ps -a -q -f ancestor=rancher/server\\)"
	utils.ExecuteRemote(appHost, command, nil, "")
}

func rancher_DockerRun(appHost string) {
	command := "docker run -d -p 8080:8080 rancher/server"
	utils.ExecuteRemote(appHost, command, nil, "")
}

func rancher_GetEnvironmentId(rancherHost string) string {
log.Println("Getting environment ID...")
	type EnvData struct {
		Data []struct {
			Id string `json:"id"`
			Name string `json:"name"`
		}
	}

	envUrl := fmt.Sprintf("http://%s/v1/accounts", rancherHost)
	log.Printf("Trying Rancher on %s...\n", envUrl)

	for {
		resp, err := http.Get(envUrl)
		if (err != nil) {
			time.Sleep(5*time.Second)
			println("Waiting for Rancher...")
			continue
		}
		envData := EnvData{}
		err = json.NewDecoder(resp.Body).Decode(&envData)

		if (err != nil && err!=io.EOF) {
			time.Sleep(5*time.Second)
			println("Waiting for Rancher...")
			continue
		}
		for _, account := range envData.Data {
			if account.Name == "Default" {
				return account.Id
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func rancher_GetApiKeys(rancherHost string, envId string) (string, string) {

	type KeyData struct {
		PublicValue string
		SecretValue string
	}
	data := map[string]interface{} {
		"type": "apikey",
		"accountId": envId,
		"name": "api_key",
		"description": "api_key",
		"created": nil,
		"kind": nil,
		"removed": nil,
		"uuid": nil,
	}
	log.Println("Getting API keys...")
	byteData, err := json.Marshal(&data)
	utils.Check(err)
	apiKeyUrl := fmt.Sprintf("http://%s/v1/projects/%s/apikey", rancherHost, envId)
	resp, err := http.Post(apiKeyUrl, "application/json", bytes.NewBuffer(byteData))
	utils.Check(err)
	keyData := KeyData{}
	err = json.NewDecoder(resp.Body).Decode(&keyData)
	utils.Check(err)
	return keyData.PublicValue, keyData.SecretValue
}

func rancher_EnableAuth(rancherHost, username, password string) {
	data := map[string]interface{}{
		"accessMode":"unrestricted",
		"name": username,
		"id": nil,
		"type":"localAuthConfig",
		"enabled": true,
		"password": password,
		"username": username,
	}

	byteData, _ := json.Marshal(&data)
	url := fmt.Sprintf("http://%s/v1/localauthconfig", rancherHost)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(byteData))
	utils.Check(err)

	println("Rancher auth enabled")
}

func rancher_RegisterRegistry(rancherHost, envId, dockerHost string) string {

	data := map[string]interface{}{
		"type": "registry",
		"serverAddress": dockerHost,
		"blockDevicePath": "",
		"created": nil,
		"description": "Private Docker Registry",
		"driverName": nil,
		"externalId": nil,
		"kind": nil,
		"name": nil,
		"removed": nil,
		"uuid": nil,
		"volumeAccessMode": nil,
	}

	byteData, _ := json.Marshal(&data)
	log.Println(string(byteData))
	type RegistryData struct{ Id string `json:"id"` }

	registryUrl := fmt.Sprintf("http://%s/v1/projects/%s/registry", rancherHost, envId)
	resp, err := http.Post(registryUrl, "application/json", bytes.NewBuffer(byteData))
	utils.Check(err)
	registryData := RegistryData{}
	err = json.NewDecoder(resp.Body).Decode(&registryData)
	utils.Check(err)
	println("Docker Registry registered")
	return registryData.Id
}

func rancher_RegistryCredentials(rancherHost, envId, registryId, username, password, email string) {

	data := map[string]interface{}{
		"type": "registryCredential",
		"registryId": registryId,
		"email": email,
		"publicValue": username,
		"secretValue": password,
		"created": nil,
		"description": nil,
		"kind": nil,
		"name": nil,
		"removed": nil,
		"uuid": nil,
	}
	credentialsUrl := fmt.Sprintf("http://%s/v1/projects/%s/registrycredential", rancherHost, envId)
	byteData, _ := json.Marshal(&data)

	_, err := http.Post(credentialsUrl, "application/json", bytes.NewBuffer(byteData))
	utils.Check(err)

	println("Docker Registry credentials configured")
}

func rancher_SetApiHost(rancherHost string) {
	type ApiData struct {
		Id    string `json:"id"`
		Links struct {
			      Self string `json:"self"`
		      } `json:"links"`
	}
	url := fmt.Sprintf("http://%s/v1/settings/api.host", rancherHost)
	resp, err := http.Get(url)
	utils.Check(err)
	apiData := ApiData{}
	err = json.NewDecoder(resp.Body).Decode(&apiData)
	utils.Check(err)

	apiUrl := apiData.Links.Self
	apiId := apiData.Id

	data := map[string]interface{}{
		"id": apiId,
		"type": "activeSetting",
		"name": "api.host",
		"activeValue": nil,
		"inDb": false,
		"source": nil,
		"value": fmt.Sprintf("http://%s", rancherHost),
	}
	byteData, _ := json.Marshal(&data)
	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(byteData))
	client := &http.Client{}
	client.Do(req)

	print("API Host set.")
}
