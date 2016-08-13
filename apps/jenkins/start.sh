#!/bin/sh

echo "2" > /var/jenkins_home/jenkins.install.UpgradeWizard.state

/bin/tini -- /usr/local/bin/jenkins.sh
