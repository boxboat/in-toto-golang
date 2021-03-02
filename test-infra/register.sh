#!/bin/bash

## DEMO REGISTRATIONS #######################################
#intotobuilder
docker exec test-infra_spire-server_1 \
/opt/spire/bin/spire-server entry create \
-selector unix:uid:1000 \
-registrationUDSPath /run/spire/sockets/spire-registration.sock \
-spiffeID spiffe://spire.boxboat.io/intoto-builder \
-parentID spiffe://spire.boxboat.io/spire/agent/sshpop/21Aic_muK032oJMhLfU1_CMNcGmfAnvESeuH5zyFw_g

## ----------------------------------------------------------------##
