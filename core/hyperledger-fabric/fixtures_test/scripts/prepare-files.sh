#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

sudo chown -R vagrant:vagrant $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data

# remove intermediatecerts
# copy admincerts and tlscerts from signcerts

# org0

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/rca-org0-admin/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/ica-org0-admin/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/admin-org0/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer0-org0/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer1-org0/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org0/orderer2-org0/msp/keystore/key.crt

# org1

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/rca-org1-admin/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/ica-org1-admin/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/admin-org1/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/user-org1/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer0-org1/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer1-org1/msp/keystore/key.crt

rm -rf $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/intermediatecerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/admincerts
cp -r $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/signcerts $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/tlscerts
mv $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/keystore/*_sk $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/data/org1/peer2-org1/msp/keystore/key.crt

# create genesis blocka and channel and set anchor peer

mkdir $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/channel-artifacts

$GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/bin/configtxgen -profile OrgsOrdererGenesis -outputBlock $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/channel-artifacts/orderer.genesis.block -channelID mycreditchain


$GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/bin/configtxgen -profile OrgsChannel -outputCreateChannelTx $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/channel-artifacts/mycreditchain.channel.tx -channelID mycreditchain


$GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/bin/configtxgen -profile OrgsChannel -outputAnchorPeersUpdate $GOPATH/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/channel-artifacts/org1.mycreditchain.anchors.tx -channelID mycreditchain -asOrg org1

