// Copyright 2016 Apcera Inc. All rights reserved.

package commands

import (
	"fmt"
	"os"

	"github.com/apcera/kurma/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	FetchCmd = &cobra.Command{
		Use:   "fetch IMAGE_URI",
		Short: "Instruct the Kurma daemon to remotely fetch and load an image",
		Run:   cmdFetch,
	}
)

func init() {
	cli.RootCmd.AddCommand(FetchCmd)
}

func cmdFetch(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	image, err := cli.GetClient().FetchImage(args[0])
	if err != nil {
		fmt.Printf("Failed to remotely fetch image %q: %s\n", args[0], err)
		os.Exit(1)
	}

	fmt.Printf("Fetched %q (%s)\n", args[0], image.Hash)
}
