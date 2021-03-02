#!/bin/bash

docker exec test-infra_spire-server_1 \
/opt/spire/bin/spire-server x509 mint \
-registrationUDSPath /run/spire/sockets/spire-registration.sock \
-spiffeID spiffe://spire.boxboat.io/$1 \
-write .

docker exec test-infra_spire-server_1 cat svid.pem > $1-svid.pem
docker exec test-infra_spire-server_1 cat key.pem > $1-key.pem
docker exec test-infra_spire-server_1 cat bundle.pem > $1-bundle.pem