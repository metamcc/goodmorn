#!/bin/bash

export FABRIC_CA_CLIENT_HOME=/data/org0/ica-org0-admin

fabric-ca-client enroll -u https://ica-org0-admin:adminpw@ica-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/ica-org0-admin

fabric-ca-client register --id.name orderer0-org0 --id.secret org0pw --id.type orderer --id.affiliation org0 --url "https://ica-org0:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/orderer0-org0

fabric-ca-client enroll -u https://orderer0-org0:org0pw@ica-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/ica-org0-admin

fabric-ca-client register --id.name orderer1-org0 --id.secret org0pw --id.type orderer --id.affiliation org0 --url "https://ica-org0:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/orderer1-org0

fabric-ca-client enroll -u https://orderer1-org0:org0pw@ica-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/ica-org0-admin

fabric-ca-client register --id.name orderer2-org0 --id.secret org0pw --id.type orderer --id.affiliation org0 --url "https://ica-org0:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/orderer2-org0

fabric-ca-client enroll -u https://orderer2-org0:org0pw@ica-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/ica-org0-admin

fabric-ca-client register --id.name admin-org0 --id.secret org0pw --id.type client --id.affiliation org0 --id.attrs 'admin=true:ecert' --url "https://ica-org0:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org0/admin-org0

fabric-ca-client enroll -u https://admin-org0:org0pw@ica-org0:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

