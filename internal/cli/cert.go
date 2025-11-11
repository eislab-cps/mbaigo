package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	certOutputDir  string
	certCommonName string
)

func init() {
	certCmd.AddCommand(certGenerateKeyCmd)
	certCmd.AddCommand(certRequestCmd)

	certGenerateKeyCmd.Flags().StringVarP(&certOutputDir, "output", "o", ".", "Output directory")

	certRequestCmd.Flags().StringVarP(&certOutputDir, "output", "o", ".", "Output directory")
	certRequestCmd.Flags().StringVarP(&certCommonName, "common-name", "n", "", "Common name (defaults to system name)")
}

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Certificate management commands",
	Long:  "Commands for generating keys and certificate signing requests",
}

var certGenerateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generate ECDSA private key",
	Long:  "Generate an ECDSA P256 private key for the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		// Marshal private key
		keyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return fmt.Errorf("failed to marshal key: %w", err)
		}

		keyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		})

		keyPath := fmt.Sprintf("%s/private_key.pem", certOutputDir)
		if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
			return fmt.Errorf("failed to write key: %w", err)
		}

		if JSON {
			return printJSON(map[string]string{"keyPath": keyPath})
		}

		printSuccess(fmt.Sprintf("Generated private key: %s", keyPath))
		fmt.Println("⚠️  Keep this file secure and do not share it!")

		return nil
	},
}

var certRequestCmd = &cobra.Command{
	Use:   "create-csr",
	Short: "Create certificate signing request",
	Long:  "Create a CSR (Certificate Signing Request) for the Certificate Authority",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig()
		if err != nil {
			return err
		}

		if certCommonName == "" {
			if name, ok := config["systemname"].(string); ok {
				certCommonName = name
			} else {
				return fmt.Errorf("common name required (use --common-name or set systemname in config)")
			}
		}

		// Load or generate private key
		keyPath := fmt.Sprintf("%s/private_key.pem", certOutputDir)
		keyPEM, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("private key not found (generate with: mbaigo cert generate-key)")
		}

		block, _ := pem.Decode(keyPEM)
		if block == nil {
			return fmt.Errorf("failed to parse PEM block")
		}

		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		// Create CSR template
		template := x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:   certCommonName,
				Organization: []string{"Arrowhead"},
			},
		}

		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
		if err != nil {
			return fmt.Errorf("failed to create CSR: %w", err)
		}

		csrPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrBytes,
		})

		csrPath := fmt.Sprintf("%s/csr.pem", certOutputDir)
		if err := os.WriteFile(csrPath, csrPEM, 0644); err != nil {
			return fmt.Errorf("failed to write CSR: %w", err)
		}

		if JSON {
			return printJSON(map[string]string{"csrPath": csrPath, "commonName": certCommonName})
		}

		printSuccess(fmt.Sprintf("Created CSR: %s", csrPath))
		fmt.Printf("Common Name: %s\n", certCommonName)
		fmt.Println("\n💡 Submit this CSR to your Certificate Authority to obtain a signed certificate")

		return nil
	},
}
