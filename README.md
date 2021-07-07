# in-toto-spiffe Proof Of Concept -- Not for Prod

in-toto is a specification to provide policy and attestations for software supply chains.
Managing the PKI around in-toto has been a major reason cited as lack of adoption.
The project aims to provide a non-production proof of concept on integrating SPIRE PKI orchestration with in-toto build chain attestation.

The integration effort required support for CA based validation of functionaries.
In-toto currently requires the corresponding public key for each private key used in the build attestation process.
This schema does not fit within most organization PKI policy.
Validation of signatures based on certificate constraints will allow end-users to integrate in-toto with existing enterprise PKI.

## Running the Demo

To run the demo, pull down the source code, install Go, and run `make test-verify`.
This will use openssl to generate a certificate chain.

SPIFFE compliant Leaf certificates are generated with SVIDs corresponding to functionaries.
These certificates are consumed by in-toto to sign link-meta data and the layout policy.

During the in-toto verification process, `certificate constraints` are checked to ensure the build step link meta-data was signed with the correct SVID.

## Building

Download the source, run `make build`.

## CLI

```text
Usage:
  in-toto [command]

Available Commands:
  help        Help about any command
  record      Creates a signed link metadata file in two steps, in order to provide
evidence for supply chain steps that cannot be carried out by a single command
  run         Executes the passed command and records paths and hashes of 'materials'
  sign        Provides command line interface to sign in-toto link or layout metadata
  verify      Verify that the software supply chain of the delivered product

Flags:
  -h, --help                              help for in-toto
      --spiffe-workload-api-path string   uds path for spiffe workload api

Use "in-toto [command] --help" for more information about a command.
```

## Example

A very simple example, just to help you starting:

```go
package main

import (
	"time"
	toto "github.com/in-toto/in-toto-golang/in_toto"
)

func main() {
	t := time.Now()
	t = t.Add(30 * 24 * time.Hour)

	var keys = make(map[string]toto.Key)

	var metablock = toto.Metablock{
		Signed: toto.Layout{
			Type: "layout",
			Expires:  t.Format("2006-01-02T15:04:05Z"),
			Steps: []toto.Step{},
			Inspect: []toto.Inspection{},
			Keys:  keys,
		},
	}

	var key toto.Key

	key.LoadKey("keys/alice", "rsassa-pss-sha256", []string{"sha256", "sha512"})

	metablock.Sign(key)

	metablock.Dump("root.layout")
}
```

### run

```text
Executes the passed command and records paths and hashes of 'materials' (i.e.
files before command execution) and 'products' (i.e. files after command
execution) and stores them together with other information (executed command,
return value, stdout, stderr, ...) to a link metadata file, which is signed
with the passed key.  Returns nonzero value on failure and zero otherwise.

Usage:
  in-toto run [flags]

Flags:
  -c, --cert string                Path to a PEM formatted certificate that corresponds with
                                   the provided key.
  -h, --help                       help for run
  -k, --key string                 Path to a PEM formatted private key file used to sign
                                   the resulting link metadata. (passing one of '--key'
                                   or '--gpg' is required) 
      --lstrip-paths stringArray   path prefixes used to left-strip artifact paths before storing
                                   them to the resulting link metadata. If multiple prefixes
                                   are specified, only a single prefix can match the path of
                                   any artifact and that is then left-stripped. All prefixes
                                   are checked to ensure none of them are a left substring
                                   of another.
  -m, --materials stringArray      Paths to files or directories, whose paths and hashes
                                   are stored in the resulting link metadata before the
                                   command is executed. Symlinks are followed.
  -n, --name string                Name used to associate the resulting link metadata
                                   with the corresponding step defined in an in-toto
                                   layout.
  -d, --output-directory string    directory to store link metadata (default "./")
  -p, --products stringArray       Paths to files or directories, whose paths and hashes
                                   are stored in the resulting link metadata after the
                                   command is executed. Symlinks are followed.

Global Flags:
      --spiffe-workload-api-path string   uds path for spiffe workload api
```

### sign

```text
Provides command line interface to sign in-toto link or layout metadata

Usage:
  in-toto sign [flags]

Flags:
  -f, --file string     Path to link or layout file to be signed or verified.
  -h, --help            help for sign
  -k, --key string      Path to PEM formatted private key used to sign the passed 
                        root layout's signature(s). Passing exactly one key using
                        '--layout-key' is required.
  -o, --output string   Path to store metadata file to be signed

Global Flags:
      --spiffe-workload-api-path string   uds path for spiffe workload api
```

### verify

```text
in-toto-verify is the main verification tool of the suite, and 
it is used to verify that the software supply chain of the delivered 
product was carried out as defined in the passed in-toto supply chain 
layout. Evidence for supply chain steps must be available in the form 
of link metadata files named ‘<step name>.<functionary keyid prefix>.link’.

Usage:
  in-toto verify [flags]

Flags:
  -h, --help                         help for verify
  -i, --intermediate-certs strings   Path(s) to PEM formatted certificates, used as intermediaries to verify
                                     the chain of trust to the layout's trusted root. These will be used in
                                     addition to any intermediates in the layout.
  -l, --layout string                Path to root layout specifying the software supply chain to be verified
  -k, --layout-keys strings          Path(s) to PEM formatted public key(s), used to verify the passed 
                                     root layout's signature(s). Passing at least one key using
                                     '--layout-keys' is required. For each passed key the layout
                                     must carry a valid signature.
  -d, --link-dir string              Path to directory where link metadata files for steps defined in 
                                     the root layout should be loaded from. If not passed links are 
                                     loaded from the current working directory.

Global Flags:
      --spiffe-workload-api-path string   uds path for spiffe workload api
```

### record

```text
Creates a signed link metadata file in two steps, in order to provide
evidence for supply chain steps that cannot be carried out by a single command
(for which ‘in-toto-run’ should be used). It returns a non-zero value on
failure and zero otherwise.

Usage:
  in-toto record [command]

Available Commands:
  start       Creates a preliminary link file recording the paths and hashes of the passed materials and signs it with the passed functionary’s key.
  stop        Records and adds the paths and hashes of the passed products to the link metadata file and updates the signature.

Flags:
  -c, --cert string   Path to a PEM formatted certificate that corresponds with the provided key.
  -h, --help          help for record
  -k, --key string    Path to a private key file to sign the resulting link metadata.
                      The keyid prefix is used as an infix for the link metadata filename,
                      i.e. ‘<name>.<keyid prefix>.link’. See ‘–key-type’ for available
                      formats. Passing one of ‘–key’ or ‘–gpg’ is required.
  -n, --name string   name for the resulting link metadata file.
                      It is also used to associate the link with a step defined
                      in an in-toto layout.

Global Flags:
      --spiffe-workload-api-path string   uds path for spiffe workload api

Use "in-toto record [command] --help" for more information about a command.
```

## Layout Certificate Constraints

Currently only URIs and common name constraints supported:

```json
{
  "cert_constraints": [{
    "uris": ["spiffe://example.com/Something"],
    "common_name": "*"
  }, {
    "uris": [],
    "common_names": ["Some User"]
  }]
}
```

## Certificate Authority

The CA for the signing keys must be included in the layout.
See [example layout](#example-layout).

## Example Layout

```json
{
 "signatures": [
 ],
 "signed": {
  "_type": "layout",
  "expires": "2021-04-03T00:00:00Z",
  "inspect": [],
  "intermediatecas": [],
  "keys": {},
  "readme": "",
  "rootcas": ["-----BEGIN CERTIFICATE-----\nMIIBkTCCATegAwIBAgIBADAKBggqhkjOPQQDAjAdMQswCQYDVQQGEwJVUzEOMAwG\nA1UEChMFU1BJUkUwHhcNMjEwMzAzMTUwMjQzWhcNMjEwNDAyMTUwMjUzWjAdMQsw\nCQYDVQQGEwJVUzEOMAwGA1UEChMFU1BJUkUwWTATBgcqhkjOPQIBBggqhkjOPQMB\nBwNCAAQ3L4PJvxT4hflEMcEcsOuvyvnOkXCH+Z5gCtDW0j6EOIBSCnFvbCf60xdF\n3jfIbVV0OVCGPRQ7QwRd5kP8vM4Jo2gwZjAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0T\nAQH/BAUwAwEB/zAdBgNVHQ4EFgQUhH+7do7BgZFg5oNTqTQhVmnfzG0wJAYDVR0R\nBB0wG4YZc3BpZmZlOi8vc3BpcmUuYm94Ym9hdC5pbzAKBggqhkjOPQQDAgNIADBF\nAiAsmHvUqQnni2OijlyCl/XONrY9C+PRjpZrVfYguBenTwIhAKsLAJHHn5MDJV+E\nYzx35oSRRRGTyM3yreDoB9G/JOPi\n-----END CERTIFICATE-----\n"],
  "steps": [
   {
    "_type": "step",
    "cert_constraints": [
      {
        "common_name": "*",
        "uris": [
          "spiffe://spire.boxboat.io/intoto-builder"
        ]
      }
    ],
    "expected_command": [
     "git",
     "clone",
     "https://gitlab.com/boxboat/demos/intoto-spire/go-hello-world"
    ],
    "expected_materials": [
     [
      "DISALLOW",
      "*"
     ]
    ],
    "expected_products": [
     [
      "CREATE",
      "*"
     ]
    ],
    "name": "clone",
    "pubkeys": [],
    "threshold": 1
   },
   {
    "_type": "step",
    "cert_constraints": [
      {
        "common_name": "*",
        "uris": [
          "spiffe://spire.boxboat.io/intoto-builder"
        ]
      }
    ],
    "expected_command": [
      "/bin/sh",
      "-c",
      "trivy --exit-code 0 --no-progress --output ./trivy-scanning-report.json --input ./go-hello-world.tar --format json"
    ],
    "expected_materials": [
     [
      "MATCH",
      "*",
      "WITH",
      "PRODUCTS",
      "FROM",
      "build-image"
     ]
    ],
    "expected_products": [
     [
      "CREATE",
      "trivy-scanning-report.json"
     ]
    ],
    "name": "scan-image",
    "pubkeys": [],
    "threshold": 1
   },
   {
    "_type": "step",
    "cert_constraints": [
      {
        "common_name": "*",
        "uris": [
          "spiffe://spire.boxboat.io/intoto-builder"
        ]
      }
    ],
    "expected_command": [
     "go",
     "build",
     "./..."
    ],
    "expected_materials": [
     [
      "MATCH",
      "*",
      "WITH",
      "PRODUCTS",
      "FROM",
      "clone"
     ],
     [
      "DISALLOW",
      "*"
     ]
    ],
    "expected_products": [
     [
      "CREATE",
      "go-hello-world"
     ],
     [
      "DISALLOW",
      "*"
     ]
    ],
    "name": "build",
    "pubkeys": [],
    "threshold": 1
   },
   {
    "_type": "step",
    "cert_constraints": [
      {
        "common_name": "*",
        "uris": [
          "spiffe://spire.boxboat.io/intoto-builder"
        ]
      }
    ],
    "expected_command": ["/bin/sh", "-c", "docker", "build", ".", "-t", "registry.gitlab.com/boxboat/demos/intoto-spire/go-hello-world", "--iidfile", "image-id", "&&", "docker", "save", "--output", "go-hello-world.tar", "registry.gitlab.com/boxboat/demos/intoto-spire/go-hello-world"],
    "expected_materials": [
     [
      "MATCH",
      "*",
      "WITH",
      "PRODUCTS",
      "FROM",
      "clone"
     ],
     [
      "DISALLOW",
      "*"
     ]
    ],
    "expected_products": [
     [
      "CREATE",
      "image-id"
     ],
     [
      "CREATE",
      "go-hello-world.tar"
     ],
     [
      "DISALLOW",
      "*"
     ]
    ],
    "name": "build-image",
    "pubkeys": [],
    "threshold": 1
   }
  ]
 }
}
```
