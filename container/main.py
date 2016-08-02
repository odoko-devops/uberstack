#/usr/bin/python

import os
import sys
import commands
import shlex
import yaml
import re
import pexpect
import rancher_server

PATH="/usr/local/bin:/bin:/usr/bin"

def step(msg):
  print 
  print "#"*80
  print msg
  print


def read_config(config_file, state_file):
  config = yaml.load(open(config_file).read())
  if os.path.exists(state_file):
    state = yaml.load(open(state_file).read())
    if state is not None:
      config.update(state)
  return config


def write_state_file(state_file, access_key=None, secret_key=None):
  if os.path.exists(state_file):
    state = yaml.load(open(state_file).read())
    if state is None:
      state = {}
  else:
    state = {}
  if access_key:
    state["rancher"] = state.get("rancher", {})
    state["rancher"]["api-access-key"] = str(access_key)
  if secret_key:
    state["rancher"] = state.get("rancher", {})
    state["rancher"]["api-secret-key"] = str(secret_key)

  with open(state_file, "w") as f:
    f.write(yaml.dump(state))


def execute(cmd, cwd=".", env={}):
  env["PATH"]=PATH
  parts=shlex.split(cmd)
  print "Executing %s %s" % (parts[0], parts[1:])
  proc = pexpect.spawn(parts[0], parts[1:], cwd=cwd, env=env)
  lines=[]
  while True:
    index = proc.expect(["(.*)\n", pexpect.TIMEOUT, pexpect.EOF], timeout=-1)
    if index == 0: 
      line = proc.match.group(1)
      print line
      lines.append(line.replace("\r", ""))
    elif index == 1:
      pass
    elif index == 2:
      break

  return "\n".join(lines)

# For some reason, execute() hangs on a call to docker-compose. This method 
# is old fashioned, but works. Note, it puts its env data into the current
# Python processes environment, which isn't great, but given docker-compose
# is the last app executed, it works.
def execute2(cmd, env):
  for key, value in env.items():
    os.environ[key]=value
  os.system(cmd)


def apply_terraform(config):
  step("Create AWS VPC Environment")

  cwd="terraform/aws"
  env = {"TF_VAR_aws_access_key": config["aws"]["access-key"],
         "TF_VAR_aws_secret_key": config["aws"]["secret-key"]}

  execute("terraform apply -state=/state/terraform.tfstate", cwd=cwd, env=env)
  vpc_id = execute("terraform output -state=/state/terraform.tfstate vpc_id", cwd=cwd, env=env)
  subnet_id = execute("terraform output -state=/state/terraform.tfstate subnet_id", cwd=cwd, env=env)
  return vpc_id, subnet_id


def destroy_terraform(config):
  step("Destroy AWS VPC Environment")

  cwd="terraform/aws"
  env = {"TF_VAR_aws_access_key": config["aws"]["access-key"],
         "TF_VAR_aws_secret_key": config["aws"]["secret-key"]}

  execute("terraform destroy -state=/state/terraform.tfstate -force", cwd=cwd, env=env)


def create_management_host(config):
  step("Create Management Docker Host")
  aws = config["aws"]
  mgt_host = aws["management-host"]
  execute('''docker-machine create --driver amazonec2
           --amazonec2-access-key=%s \
           --amazonec2-secret-key=%s \
               --amazonec2-vpc-id=%s \
               --amazonec2-instance-type %s \
               --amazonec2-security-group management-tools \
               --amazonec2-region %s \
               --amazonec2-zone %s \
               --amazonec2-subnet-id %s \
               --amazonec2-tags name=management-tools \
               --amazonec2-ssh-keypath %s \
           management''' % (aws["access-key"],
                            aws["secret-key"],
                vpc_id, 
                            mgt_host["instance-type"], 
                aws["region"], 
                aws["zone"], 
                subnet_id, 
                "/id_rsa"))
  instance_id=execute("docker-machine inspect management -f '{{.Driver.InstanceId}}'")
  return instance_id


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


def create_docker_host_with_docker_machine(config, count):
  step("Create Docker Host %s" % count )
  aws = config["aws"]
  host = aws["docker-host"]
  execute('''docker-machine create
                  --driver amazonec2 \
                  --amazonec2-access-key %s \
                  --amazonec2-secret-key %s \
                      --amazonec2-vpc-id %s \
                      --amazonec2-instance-type %s \
                      --amazonec2-security-group management-tools \
                      --amazonec2-region %s \
                      --amazonec2-zone %s \
                      --amazonec2-subnet-id %s \
                      --amazonec2-tags name=management-tools \
                      --amazonec2-ssh-keypath %s \
                      docker-host%s
            ''' % (aws["access-key"],
                   aws["secret-key"],
                   vpc_id,
                   host["instance-type"],
                   aws["region"],
                   aws["zone"],
                   subnet_id,
                   "/id_rsa",
                   count))

  rancher = config["rancher"]
  execute("docker-machine scp install-rancher-agent.sh docker-host%s:" % count)
  execute("docker-machine ssh docker-host%s ./install-rancher-agent.sh %s %s %s" %
          (count,
           config["apps"]["rancher"]["name"],
           rancher["api-access-key"],
           rancher["api-secret-key"]))


def create_docker_host_with_rancher_cli(config):
  step("Create Management Docker Host")
  aws = config["aws"]
  host = aws["docker-host"]
  rancher = config["rancher"]
  execute('''rancher --url http://%s/v1 \
                     --access-key %s \
                 --secret-key %s \
                 host create \
                 --driver amazonec2 \
                 --amazonec2-access-key %s \
                 --amazonec2-secret-key %s \
                     --amazonec2-vpc-id %s \
                     --amazonec2-instance-type %s \
                     --amazonec2-security-group management-tools \
                     --amazonec2-region %s \
                     --amazonec2-zone %s \
                     --amazonec2-subnet-id %s \
                     --amazonec2-tags name=management-tools \
                     --amazonec2-ssh-keypath %s \
           ''' % (config["apps"]["rancher"]["name"],
                  rancher["api-access-key"],
                  rancher["api-secret-key"],
              aws["access-key"],
                  aws["secret-key"],
              vpc_id, 
                  host["instance-type"], 
              aws["region"], 
              aws["zone"], 
              subnet_id, 
              "/id_rsa"))


def make_elastic_ip_association(config, instance_id, eip_allocation):
  step("Associate predefined EIP with Docker Host")
  env = {"TF_VAR_aws_access_key": config["aws"]["access-key"],
         "TF_VAR_aws_secret_key": config["aws"]["secret-key"],
         "TF_VAR_instance_id": instance_id, 
     "TF_VAR_allocation_id": eip_allocation}
  execute("terraform apply", cwd="terraform/aws-eip", env=env)


def docker_machine_destroy(host):
  step("Destroy host: %s" % host)
  execute('docker-machine rm -f %s' % host)


def add_ubuntu_to_docker_group(host):
  step("Add ubuntu user to docker unix group on host %s")
  execute('docker-machine ssh %s "sudo gpasswd -a ubuntu docker"' % host)


def create_jenkins_mount_point():
  step("Create Mount Point for Jenkins")
  execute('docker-machine ssh management "sudo mkdir /jenkins ; sudo chown 1000 /jenkins"')


def get_docker_environment(host):
  RE=re.compile(r"export (.*)=\"(.*)\"")
  execute("docker-machine regenerate-certs -f management")
  result=execute("docker-machine env --shell management")

  env={}
  for line in result.split("\n"):
    m=RE.match(line)
    if m:
      env[m.group(1)] = m.group[2]
  return env


def docker_compose(config):
  step("Deploy Management Services")
  env={
    "DOCKER_HOSTNAME": config["apps"]["docker"]["name"],
    "RANCHER_HOSTNAME": config["apps"]["rancher"]["name"],
    "JENKINS_HOSTNAME": config["apps"]["jenkins"]["name"],
    "EMAIL": config["auth"]["email"],
    "USERNAME": config["auth"]["username"],
    "PASSWORD": config["auth"]["password"],
    "DOCKER_TLS_VERIFY": "1",
    "DOCKER_HOST": "tcp://%s:2376" % config["aws"]["management-host"]["elastic-ip"],
    "DOCKER_CERT_PATH": "/odoko/.docker/machine/machines/management",
    "DOCKER_MACHINE_NAME": "management"
  }
  #get_docker_environment("management")
  execute2("docker-compose up -d", env=env)



if __name__ == "__main__":
  config_file = "/config.yml"
  state_file = "/state/state.yml"
  config = read_config(config_file, state_file)

  action = sys.argv[1]

  if action == "destroy":
    docker_machine_destroy("management")
    destroy_terraform(config)

  elif action == "up":
    vpc_id, subnet_id = apply_terraform(config)

    instance_id = create_management_host(config)
    make_elastic_ip_association(config, instance_id, config["aws"]["management-host"]["elastic-ip-allocation"])
    add_ubuntu_to_docker_group("management")
    create_jenkins_mount_point()
    execute("docker-machine regenerate-certs -f management")
    docker_compose(config)
    access_key, secret_key = configure_rancher(config)
    write_state_file(state_file, access_key=access_key, secret_key=secret_key)

  elif action == "docker-up":
    vpc_id, subnet_id = apply_terraform(config)
    #create_docker_host_with_rancher_cli(config)
    count = int(config["hosts"]["count"])
    print "Creating %s hosts" % count
    for i in range(0, count):
      create_docker_host_with_docker_machine(config, count)
  else:
    print "Unknown action: %s" % action
