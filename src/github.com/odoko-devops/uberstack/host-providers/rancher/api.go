package rancher

import (
	"time"
	"io"
	"encoding/json"
	"net/http"
	"log"
	"fmt"
)

func getEnvironmentId(rancherHost, accessKey, secretKey string, environmentName string) string {
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
		client := &http.Client{}
		req, err := http.NewRequest("GET", envUrl, nil)
		req.SetBasicAuth(accessKey, secretKey)
		resp, err := client.Do(req)
		if (err != nil) {
			time.Sleep(5*time.Second)
			println("Waiting for Rancher...")
			continue
		}
		envData := EnvData{}
		err = json.NewDecoder(resp.Body).Decode(&envData)

		if (err != nil && err!=io.EOF) {
			time.Sleep(5*time.Second)
			log.Println("Waiting for Rancher...")
			continue
		}
		for _, account := range envData.Data {
			if account.Name == environmentName {
				return account.Id
			}
		}
		log.Printf("No account found. Retrying (%s, %s)", envData, resp.Body)
		time.Sleep(5 * time.Second)
	}
}