#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# Initialize the intermediate CA
fabric-ca-server init -b $BOOTSTRAP_USER_PASS -u $PARENT_URL

# Copy the intermediate CA's certificate chain to the data directory to be used by others
cp $FABRIC_CA_SERVER_HOME/ca-chain.pem $TARGET_CHAINFILE

# Add the custom orgs
#for o in $FABRIC_ORGS; do
#   aff=$aff"\n   $o: []"
#done
#aff="${aff#\\n   }"
#sed -i "/affiliations:/a \\   $aff" \
#   $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml

sed -i "s/  type: sqlite3/  type: mysql/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/  datasource: fabric-ca-server.db/  datasource: root:qawsedrf@tcp(icamariadb-org1:3306)\/fabric_ca?parseTime=true/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "/   org1:/,+2 d" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/   org2:/   org0: []\n   org1: []\n   org1:/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sleep 10s

# Start the intermediate CA
fabric-ca-server start
