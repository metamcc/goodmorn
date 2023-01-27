#!/bin/bash

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client enroll -u https://ica-org1-admin:adminpw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client register --id.name peer0-org1 --id.secret org1pw --id.type peer --id.affiliation org1 --url "https://ica-org1:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/peer0-org1

fabric-ca-client enroll -u https://peer0-org1:org1pw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client register --id.name peer1-org1 --id.secret org1pw --id.type peer --id.affiliation org1 --url "https://ica-org1:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/peer1-org1

fabric-ca-client enroll -u https://peer1-org1:org1pw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client register --id.name peer2-org1 --id.secret org1pw --id.type peer --id.affiliation org1 --url "https://ica-org1:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/peer2-org1

fabric-ca-client enroll -u https://peer2-org1:org1pw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client register --id.name admin-org1 --id.secret org1pw --id.type client --id.affiliation org1 --id.attrs 'hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert' --url "https://ica-org1:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/admin-org1

fabric-ca-client enroll -u https://admin-org1:org1pw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/ica-org1-admin

fabric-ca-client register --id.name user-org1 --id.secret org1pw --id.type client --id.affiliation org1 --url "https://ica-org1:7054" --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem

export FABRIC_CA_CLIENT_HOME=/data/org1/user-org1

fabric-ca-client enroll -u https://user-org1:org1pw@ica-org1:7054 --tls.certfiles $FABRIC_CA_SERVER_HOME/ca-cert.pem
