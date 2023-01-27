#!/bin/bash

export FABRIC_CA_CLIENT_HOME=/data/org0/rca-org0-admin

fabric-ca-client enroll -u https://rca-org0-admin:adminpw@rca-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem
