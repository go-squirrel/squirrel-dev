package app

import (
	"github.com/spf13/cobra"

	"squirrel-dev/internal/squctl/certs"
)

// NewCertsCommand 创建 certs 子命令
func NewCertsCommand() *cobra.Command {
	o := certs.NewCertsOptions()

	cmd := &cobra.Command{
		Use:   "certs",
		Short: "Generate certificates for mTLS authentication",
		Long: `Generate CA, server, and client certificates for mTLS authentication.

Examples:
  # Generate all certificates with default settings
  squctl certs

  # Generate only CA certificate
  squctl certs --only-ca

  # Generate certificates with custom output directory
  squctl certs --output /etc/squirrel/certs

  # Generate server certificate with specific hosts
  squctl certs --server-hosts 127.0.0.1,localhost,apiserver.example.com

  # Generate client certificate with custom CN
  squctl certs --client-cn squirrel-agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	// 输出目录
	cmd.Flags().StringVarP(&o.OutputDir, "output", "o", "./certs", "Output directory for certificates")

	// CA 相关
	cmd.Flags().BoolVar(&o.OnlyCA, "only-ca", false, "Generate only CA certificate")
	cmd.Flags().StringVar(&o.CACN, "ca-cn", "squirrel-ca", "Common Name for CA certificate")

	// Server 证书相关
	cmd.Flags().BoolVar(&o.OnlyServer, "only-server", false, "Generate only server certificate (requires existing CA)")
	cmd.Flags().StringVar(&o.ServerCN, "server-cn", "squirrel-apiserver", "Common Name for server certificate")
	cmd.Flags().StringSliceVar(&o.ServerHosts, "server-hosts", []string{"127.0.0.1", "localhost"}, "Hosts (IP/DNS) for server certificate, comma separated")

	// Client 证书相关
	cmd.Flags().BoolVar(&o.OnlyClient, "only-client", false, "Generate only client certificate (requires existing CA)")
	cmd.Flags().StringVar(&o.ClientCN, "client-cn", "squirrel-agent", "Common Name for client certificate")

	// 通用参数
	cmd.Flags().DurationVar(&o.Expiry, "expiry", certs.DefaultExpiry, "Certificate validity period (e.g. 87600h for 10 years)")
	cmd.Flags().IntVar(&o.KeySize, "key-size", 2048, "RSA key size in bits")
	cmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "Overwrite existing certificates")

	return cmd
}
