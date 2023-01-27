#!/bin/bash

export FABRIC_CA_CLIENT_HOME=/data/org1/rca-org1-admin

fabric-ca-client enroll -u https://rca-org1-admin:adminpw@rca-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem
