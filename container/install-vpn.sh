#!/bin/sh

echo CIDR=${CIDR}
echo PUBLIC_IP=${PUBLIC_IP}
echo USERNAME=${USERNAME}

docker-machine ssh management "sudo mkdir -p /etc/openvpn"
echo "Container for data:"
docker run --name ovpn-data -v /etc/openvpn:/etc/openvpn busybox
echo "Generate Configs for ${CIDR} and ${PUBLIC_IP}"
docker run --volumes-from ovpn-data --rm gosuri/openvpn ovpn_genconfig -p ${CIDR} -u udp://${PUBLIC_IP}
echo "Init the PKI"
docker run --volumes-from ovpn-data --rm -it gosuri/openvpn ovpn_initpki
echo "Start Listening""
docker run --volumes-from ovpn-data -d -p 1194:1194/udp --cap-add=NET_ADMIN gosuri/openvpn
echo "Build client"
docker run --volumes-from ovpn-data --rm -it gosuri/openvpn easyrsa build-client-full "${USERNAME}" nopass
echo "Download client"
docker run --volumes-from ovpn-data --rm     gosuri/openvpn ovpn_getclient "\${USERNAME}" > /state/${USERNAME}.ovpn
