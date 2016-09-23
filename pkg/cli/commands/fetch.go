// Copyright 2016 Apcera Inc. All rights reserved.

package commands

import (
	"fmt"
	"os"

	"github.com/apcera/kurma/pkg/cli"
	"github.com/apcera/kurma/pkg/image"
	"github.com/spf13/cobra"
)

var (
	FetchCmd = &cobra.Command{
		Use:   "fetch IMAGE_URI",
		Short: "Instruct the Kurma daemon to remotely fetch and load an image",
		Run:   cmdFetch,
	}

	insecureImageFetch bool
)

func init() {
	cli.RootCmd.AddCommand(FetchCmd)
	// TODO: insecure option should not be true by default.
	FetchCmd.Flags().BoolVarP(&insecureFetch, "insecure", "", true, "pull without verifying signature or enforcing HTTPS")
}

func cmdFetch(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	image, err := cli.GetClient().FetchImage(args[0], &image.FetchConfig{Insecure: insecureImageFetch})
	if err != nil {
		fmt.Printf("Failed to remotely fetch image %q: %s\n", req.ImageURI, err)
		os.Exit(1)
	}

	fmt.Printf("Fetched %q (%s)\n", req.ImageURI, image.Hash)
}
