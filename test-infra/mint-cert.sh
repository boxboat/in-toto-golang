#!/bin/bash




docker exec test-infra_spire-server_1 \
/opt/spire/bin/spire-server x509 mint \
-registrationUDSPath /run/spire/sockets/spire-registration.sock \
-spiffeID spiffe://spire.boxboat.io/$1