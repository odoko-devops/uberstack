#!/bin/sh

cat > /tmp/rancher-data <<EOF
{
  "accessMode":"unrestricted",
  "name":"${USERNAME}",
  "id":null,
  "type":"localAuthConfig",
  "enabled":true,
  "password":"${PASSWORD}",
  "username":"${USERNAME}"}
EOF

while true; do
  curl -s http://rancher.odoko.org/v1/localauthconfig -d @/tmp/rancher-data -H "Content-type: application/json" > /dev/null
  if [ $? = 0 ]; then 
    break
  fi
  
  echo "Waiting for Rancher to start..."
  sleep 10
done

echo "Rancher auth enabled"
