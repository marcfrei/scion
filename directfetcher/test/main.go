package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/scionproto/scion/directfetcher"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/private/topology"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <topology-file> <destination-IA>", os.Args[0])
	}

	topologyFile := os.Args[1]
	dstIAStr := os.Args[2]

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

	fetcher := directfetcher.New(topo)

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
