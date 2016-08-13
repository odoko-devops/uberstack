#!groovy

import jenkins.model.*
import hudson.security.*

def instance = Jenkins.getInstance()

def username = System.getenv("USERNAME")
def password = System.getenv("PASSWORD")
def docker = System.getenv("DOCKER_HOSTNAME")

println("STARTING DOCKER LOGIN")
def cmd = "docker login -u $username -p $password $docker"
println("logging into Docker with: $cmd")
while (true) {
  def process = cmd.execute();
  process.waitFor();
  println(process.text)
  if (process.exitValue() == 0) {
    println("Logged into Docker")
    break
  }
  println("Waiting for Docker...x")
  sleep(5000)
}
