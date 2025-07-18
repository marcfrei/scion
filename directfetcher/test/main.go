package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/scionproto/scion/directfetcher"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/drkey"
	"github.com/scionproto/scion/pkg/drkey/generic"
	"github.com/scionproto/scion/private/storage"
	"github.com/scionproto/scion/private/storage/trust/fspersister"
	"github.com/scionproto/scion/private/topology"
	"github.com/scionproto/scion/private/trust"
)

const persistTRCs = false

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  showpaths  Show paths to a destination\n")
		fmt.Fprintf(os.Stderr, "  showkeys   Show DRKey keys\n")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "showpaths":
		showPaths()
	case "showkeys":
		showKeys()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func showPaths() {
	var topologyFile string
	var certsDir string
	var dstIAStr string

	showPathsCmd := flag.NewFlagSet("showpaths", flag.ExitOnError)
	showPathsCmd.StringVar(&topologyFile, "topo", "topology.json", "Path to topology file (required)")
	showPathsCmd.StringVar(&certsDir, "certs", "", "Path to certs directory (optional, enables verification)")
	showPathsCmd.StringVar(&dstIAStr, "dst", "", "Destination ISD-AS (required)")
	showPathsCmd.Parse(os.Args[2:])

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

func showKeys() {
	var topologyFile string
	var mode string
	var serverIAStr string
	var clientIAStr string
	var serverHost string
	var clientHost string
	var protoID int

	showKeysCmd := flag.NewFlagSet("showkeys", flag.ExitOnError)
	showKeysCmd.StringVar(&topologyFile, "topo", "topology.json", "Path to topology file (required)")
	showKeysCmd.StringVar(&mode, "mode", "server", "Key mode: server, client")
	showKeysCmd.StringVar(&serverIAStr, "server-ia", "", "Server ISD-AS (required)")
	showKeysCmd.StringVar(&clientIAStr, "client-ia", "", "Client ISD-AS (required)")
	showKeysCmd.StringVar(&serverHost, "server-host", "", "Server host (required)")
	showKeysCmd.StringVar(&clientHost, "client-host", "", "Client host (required)")
	showKeysCmd.IntVar(&protoID, "proto", 123, "Protocol ID")
	showKeysCmd.Parse(os.Args[2:])

	if topologyFile == "" {
		log.Fatalf("Path to topology file is required (use -topo flag)")
	}
	if serverIAStr == "" {
		log.Fatalf("Server ISD-AS is required (use -server flag)")
	}
	if clientIAStr == "" {
		log.Fatalf("Client ISD-AS is required (use -client flag)")
	}
	if serverHost == "" {
		log.Fatalf("Server host is required (use -serverhost flag)")
	}
	if clientHost == "" {
		log.Fatalf("Client host is required (use -clienthost flag)")
	}

	topo, err := topology.NewLoader(topology.LoaderCfg{
		File: topologyFile,
	})
	if err != nil {
		log.Fatalf("Failed to create topology loader: %v", err)
	}

	serverIA, err := addr.ParseIA(serverIAStr)
	if err != nil {
		log.Fatalf("Failed to parse server ISD-AS '%s': %v", serverIAStr, err)
	}

	clientIA, err := addr.ParseIA(clientIAStr)
	if err != nil {
		log.Fatalf("Failed to parse client ISD-AS '%s': %v", clientIAStr, err)
	}

	fetcher := directfetcher.New(topo, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
	defer cancel()

	validity := time.Now()

	switch mode {
	case "server":
		hostASMeta := drkey.HostASMeta{
			ProtoId:  drkey.Protocol(protoID),
			Validity: validity,
			SrcIA:    serverIA,
			DstIA:    clientIA,
			SrcHost:  serverHost,
		}
		hostASKey, err := fetcher.FetchHostASKey(ctx, hostASMeta)
		if err != nil {
			log.Fatalf("Failed to fetch host-AS key: %v", err)
		}
		t0 := time.Now()
		serverKey, err := deriveHostHostKey(hostASKey, clientHost)
		if err != nil {
			log.Fatalf("Failed to derive host-host key: %v", err)
		}
		durationServer := time.Since(t0)
		fmt.Printf("Server\thost key = %s\tduration = %s\n", hex.EncodeToString(serverKey.Key[:]), durationServer)
	case "client":
		hostHostMeta := drkey.HostHostMeta{
			ProtoId:  drkey.Protocol(protoID),
			Validity: validity,
			SrcIA:    serverIA,
			DstIA:    clientIA,
			SrcHost:  serverHost,
			DstHost:  clientHost,
		}
		t0 := time.Now()
		clientKey, err := fetcher.FetchHostHostKey(ctx, hostHostMeta)
		if err != nil {
			log.Fatalf("Failed to fetch host-host key: %v", err)
		}
		durationClient := time.Since(t0)
		fmt.Printf("Client\thost key = %s\tduration = %s\n", hex.EncodeToString(clientKey.Key[:]), durationClient)
	default:
		log.Fatalf("Unknown mode: %s (supported: server, client)", mode)
	}
}

func deriveHostHostKey(hostASKey drkey.HostASKey, dstHost string) (
	drkey.HostHostKey, error) {
	deriver := generic.Deriver{
		Proto: hostASKey.ProtoId,
	}
	hostHostKey, err := deriver.DeriveHostHost(
		dstHost,
		hostASKey.Key,
	)
	if err != nil {
		return drkey.HostHostKey{}, err
	}
	return drkey.HostHostKey{
		ProtoId: hostASKey.ProtoId,
		Epoch:   hostASKey.Epoch,
		SrcIA:   hostASKey.SrcIA,
		DstIA:   hostASKey.DstIA,
		SrcHost: hostASKey.SrcHost,
		DstHost: dstHost,
		Key:     hostHostKey,
	}, nil
}
