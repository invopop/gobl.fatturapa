package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	sdi "github.com/invopop/gobl.fatturapa/sdi"
	"github.com/spf13/cobra"
)

type transmitOpts struct {
	*rootOpts
	config *sdi.Config
}

func transmit(o *rootOpts) *transmitOpts {
	return &transmitOpts{rootOpts: o}
}

func (c *transmitOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transmit [file] [environment]",
		Short: "Transmit a FatturaPA XML file to SdI in selected environment",
		RunE:  c.runE,
	}

	f := cmd.Flags()
	f.Bool("verbose", false, "Logs all requests into the console")
	f.String("ca-cert", "", "Path to a file containing the CA certificate")
	f.String("cert", "", "Path to a file containing the SDI PEM certificate")
	f.String("key", "", "Path to a file containing the sender PEM RSA private key")
	f.StringP("env", "e", "test", "Environment for running command")
	_ = cmd.MarkFlagRequired("ca-cert")
	_ = cmd.MarkFlagRequired("cert")
	_ = cmd.MarkFlagRequired("key")

	return cmd
}

func (c *transmitOpts) runE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	config, err := getConfigByEnv(cmd)
	if err != nil {
		return err
	}

	caCertPool, err := loadDataFromFlag(cmd, "ca-cert")
	if err != nil {
		return err
	}

	certPEM, err := loadDataFromFlag(cmd, "cert")
	if err != nil {
		return err
	}

	keyPEM, err := loadDataFromFlag(cmd, "key")
	if err != nil {
		return err
	}

	client := sdi.NewClient(
		sdi.WithDebugMode(verbose),
		sdi.WithCACertPool(caCertPool),
		sdi.WithClientCertPair(certPEM, keyPEM),
	)

	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	c.config = config

	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	fileName := filepath.Base(args[0])

	invOpts := sdi.InvoiceOpts{FileName: fileName, FileBody: data}

	resp, err := sdi.SendInvoice(ctx, invOpts, client, *c.config)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	fmt.Printf("Invoice sent to SdI.\nHere is the response: %v\n", resp.Body.Response)

	return nil
}

func getConfigByEnv(cmd *cobra.Command) (*sdi.Config, error) {
	env, _ := cmd.Flags().GetString("env")
	configs := map[string]*sdi.Config{
		"dev":  &sdi.DevelopmentSdIConfig,
		"test": &sdi.TestSdIConfig,
		"prod": &sdi.ProductionSdIConfig,
	}
	config, exists := configs[env]
	if !exists {
		return nil, fmt.Errorf("invalid environment: %s. Allowed values are 'dev', 'test', 'prod'", env)
	}

	return config, nil
}

func loadDataFromFlag(cmd *cobra.Command, flagName string) ([]byte, error) {
	filePath, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
