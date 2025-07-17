package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/scionproto/scion/directfetcher"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/private/storage"
	"github.com/scionproto/scion/private/storage/trust/fspersister"
	"github.com/scionproto/scion/private/topology"
	"github.com/scionproto/scion/private/trust"
)

const persistTRCs = false

func main() {
	var topologyFile string
	var certsDir string
	var dstIAStr string

	flag.StringVar(&topologyFile, "topo", "topology.json", "Path to topology file (required)")
	flag.StringVar(&certsDir, "certs", "", "Path to certs directory (optional, enables verification)")
	flag.StringVar(&dstIAStr, "dst", "", "Destination ISD-AS (required)")
	flag.Parse()

	if topologyFile == "" {
		log.Fatalf("Path to topology file is required (use -topo flag)")
	}
	if dstIAStr == "" {
		log.Fatalf("Destination ISD-AS is required (use -dst flag)")
	}

	topo, err := topology.NewLoader(topology.LoaderCfg{
		File: topologyFile,
	})
	if err != nil {
		log.Fatalf("Failed to create topology loader: %v", err)
	}

	dstIA, err := addr.ParseIA(dstIAStr)
	if err != nil {
		log.Fatalf("Failed to parse destination ISD-AS '%s': %v", dstIAStr, err)
	}

	var trustDB storage.TrustDB
	if certsDir != "" {
		trustDB, err = storage.NewTrustStorage(storage.DBConfig{
			Connection: "file:trust.db?mode=memory&cache=shared",
		})
		if err != nil {
			log.Fatalf("Failed to create trust DB: %v", err)
		}

		if persistTRCs {
			trustDB = fspersister.WrapDB(trustDB, fspersister.Config{
				TRCDir: certsDir,
			})
		}

		_, err = trust.LoadTRCs(context.Background(), certsDir, trustDB)
		if err != nil {
			log.Fatalf("Failed to load TRCs: %v", err)
		}
	}

	fetcher := directfetcher.New(topo, trustDB)

	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()

	paths, err := fetcher.FetchPaths(ctx, dstIA)
	if err != nil {
		log.Fatalf("Failed to fetch paths: %v", err)
	}

	for i, path := range paths {
		fmt.Printf("[%d] %v\n", i, path)
	}
}
