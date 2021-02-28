package in_toto

import (
	"context"
	"log"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

//GetSVID grabs the x.509 context.
func GetSVID(ctx context.Context, socketPath string) Key {

	var k Key

	//*x509.Certificate

	client, err := workloadapi.New(ctx, workloadapi.WithAddr(socketPath))
	if err != nil {
		log.Fatalf("Unable to create workload API client: %v", err)
	}
	defer client.Close()

	svidContext, err := client.FetchX509Context(ctx)
	if err != nil {
		log.Fatalf("Error grabbing x.509 context: %v", err)
	}

	log.Printf("using svid %v", svidContext.DefaultSVID().ID.String())

	parent := svidContext.Bundles.Bundles()[0].X509Authorities()[1].URIs[0]
	log.Printf("parent svid %v", parent.String())

	svid, keyBytes, err := svidContext.DefaultSVID().Marshal()
	if err != nil {
		log.Fatalf("Error marshaling certificate: %v", err)
	}

	k.KeyVal.Private = string(keyBytes)
	k.KeyVal.Public = string(svid)
	k.KeyVal.Certificate = string(svid)
	k.Scheme = ecdsaSha2nistp384
	k.KeyType = ecdsaKeyType //this should be known from the ASN.1 data
	k.KeyIDHashAlgorithms = []string{"sha256", "sha512"}
	err = k.generateKeyID()
	if err != nil {
		log.Fatalf("Error generating keyID: %v", err)
	}

	return k
}
