package apps

var rancherDockerCompose = `
rancher:
  image: rancher/server
  ports:
    - 8080:8080
`

/*
def configure_rancher(config):
  rancher_host = config["apps"]["rancher"]["name"]
  docker_host = config["apps"]["docker"]["name"]
  email = config["auth"]["email"]
  username = config["auth"]["username"]
  password = config["auth"]["password"]

  rancher_server.wait_for_rancher(rancher_host)
  rancher_server.set_api_host(rancher_host)
  rancher_server.register_docker_registry(rancher_host, docker_host, email, username, password)
  access_key, secret_key = rancher_server.get_keys(rancher_host)
  rancher_server.enable_auth(rancher_host, username, password)

  return access_key, secret_key


import requests
import pyjq
import os
import time
import logging

##########################################################
# Logging and statics
DELAY=10
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


##########################################################
# HTTP Utility methods
def get(url, path=None):
  while True:
    try:
      logger.debug("GET URL: %s", url)
      json = requests.get(url).json()
      logger.debug("Response: %s", json)

      if path is None:
        return json
      else:
        data = pyjq.first(path, json)
        logger.debug("%s reveals %s", path, data)
        return data
    except Exception, e:
      logger.warn("%s, retrying in %ss", e, DELAY)
      time.sleep(DELAY)


def post(url, data, path=None):
  while True:
    try:
      logging.debug("POST URL: %s", url)
      logging.debug("Data: %s", data)
      json = requests.post(url, json=data).json()
      logging.debug("Response: %s", json)

      if path is None:
        return json
      else:
        data = pyjq.first(path, json)
        logging.debug("%s reveals %s", path, data)
        return data
    except Exception, e:
      logger.warn("%s, retrying in %ss", e, DELAY)
      time.sleep(DELAY)


def put(url, data):
  while True:
    try:
      logging.debug("PUT URL: %s", url)
      logging.debug("Data: %s", data)
      requests.put(url, data)
      return
    except Exception, e:
      logger.warn("%s, retrying in %ss", e, DELAY)
      time.sleep(DELAY)


##########################################################
# Rancher Utility methods

def get_rancher_env(rancher_host):
  env_url = "http://%s/v1/accounts" % rancher_host
  return get(env_url, ".data | map(select(.name == \"Default\")) | .[0].id")


##########################################################
# Actually do it

def wait_for_rancher(rancher_host):
  url = "http://%s/v1/settings" % rancher_host
  while True:
    try:
      data = requests.get(url)
      break
    except:
      print "Waiting for Rancher..."
      time.sleep(5)
  print "Rancher detected."


def get_api_url(rancher_host):
  url = "http://%s/v1/settings/api.host" % rancher_host
  data = get(url, '{link: .links.self, id: .id}')
  logging.debug("API DATA: %s", data)
  return data["link"], data["id"]


def set_api_host(rancher_host):
  api_url,id = get_api_url(rancher_host)
  data={
  "id": id, # such as: "1as!api.host",
  "type": "activeSetting",
  "name": "api.host",
  "activeValue": None,
  "inDb": False,
  "source": None,
  "value": "http://%s" % rancher_host
  }

  put(api_url, data)

  print "API Host set."


def register_docker_registry(rancher_host, docker_host, email, username, password):
  env = get_rancher_env(rancher_host)
  data={
  "type": "registry",
  "serverAddress": docker_host,
  "blockDevicePath": None,
  "created": None,
  "description": "Private Docker Registry",
  "driverName": None,
  "externalId": None,
  "kind": None,
  "name": None,
  "removed": None,
  "uuid": None,
  "volumeAccessMode": None
  }

  registry_url = "http://%s/v1/projects/%s/registry" % (rancher_host, env)
  registry_id = post(registry_url, data, ".id")

  data={
  "type": "registryCredential",
  "registryId": registry_id,
  "email": email,
  "publicValue": username,
  "secretValue": password,
  "created": None,
  "description": None,
  "kind": None,
  "name": None,
  "removed": None,
  "uuid": None
  }
  credentials_url = "http://%s/v1/projects/%s/registrycredential" % (rancher_host, env)
  post(credentials_url, data)

  print "Docker registry registered"


def get_keys(rancher_host):
  env = get_rancher_env(rancher_host)

  data={
      "type":"apikey",
      "accountId": env,
      "name":"api_key",
      "description": "api_key",
      "created": None,
      "kind": None,
      "removed": None,
      "uuid":None
  }
  key_url = "http://%s/v1/projects/%s/apikey" % (rancher_host, env)
  result = post(key_url, data, "{access_key: .publicValue, secret_key: .secretValue}")
  return result["access_key"], result["secret_key"]


def enable_auth(rancher_host, username, password):
  data={
  "accessMode":"unrestricted",
  "name": username,
  "id":None,
  "type":"localAuthConfig",
  "enabled": True,
  "password": password,
  "username": username
  }

  url = "http://%s/v1/localauthconfig" % rancher_host
  post(url, data)

  print "Rancher auth enabled"


*/
