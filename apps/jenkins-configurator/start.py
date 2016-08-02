#!/usr/bin/env python

import requests
import os
import time
import logging

##########################################################
# Logging
logging.basicConfig(level=logging.DEBUG)


##########################################################
# HTTP Utility methods
def get(url):
  logging.debug("GET URL: %s", url)
  text = requests.get(url).text
  logging.debug("Response: %s", text)


def post_xml(url, data):
  logging.debug("POST URL: %s", url)
  logging.debug("Data: %s", data)
  text = requests.post(url, data=data, headers={"Content-type": "text/xml"}).text
  logging.debug("Response: %s", text)
  

##########################################################
# Install Git Plugin
def install_git_plugin(jenkins_host):
  url = "http://%s/pluginManager/installNecessaryPlugins" % jenkins_host
  data='<jenkins><install plugin="git@2.0" /></jenkins>'
  post_xml(url, data)
  print "Git plugin installed"


##########################################################
# Log into Docker
def log_into_docker(jenkins_host, docker_host, username, password):
  data="""<?xml version='1.0' encoding='UTF-8'?>
  <project>
    <actions/>
    <description></description>
    <keepDependencies>false</keepDependencies>
    <properties/>
    <scm class="hudson.scm.NullSCM"/>
    <canRoam>true</canRoam>
    <disabled>false</disabled>
    <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>
    <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>
    <triggers/>
    <concurrentBuild>false</concurrentBuild>
    <builders>
      <hudson.tasks.Shell>
        <command>docker login -u $USERNAME -p %s$PASSWORD $DOCKER_HOSTNAME</command>
      </hudson.tasks.Shell>
    </builders>
    <publishers/>
    <buildWrappers/>
  </project>
  """

  # Create Job
  url = "http://%s/createItem?name=docker-login" % jenkins_host
  post_xml(url, data)
  # Actually login
  get("http://%s/job/docker-login/buildWithParameters?USERNAME=%s&PASSWORD=%s&DOCKER_HOSTNAME=%s" % (username, password, docker_host))
  # Delete job
  get("http://%s/job/docker-login/doDelete" % jenkins_host)
  print "Logged into Docker."

if __name__ == "__main__":
  jenkins_host = os.getenv("JENKINS")
  docker_host = os.getenv("DOCKER_HOSTNAME")
  username = os.getenv("USERNAME")
  password = os.getenv("PASSWORD")

  install_git_plugin(jenkins_host)
  #log_into_docker(jenkins_host, docker_host, username, password)
  print "Done."
