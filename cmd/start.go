package cmd

import (
	"fmt"
	"github.com/mberwanger/demux-proxy/internal/context"
	"github.com/mberwanger/demux-proxy/internal/proxy"
	"github.com/spf13/cobra"
)

type startCmd struct {
	cmd  *cobra.Command
	opts startOpts
}

type startOpts struct {
	bindAddress string
}

func newStartCmd() *startCmd {
	var root = &startCmd{}
	var cmd = &cobra.Command{
		Use:           "start",
		Aliases:       []string{"c"},
		Short:         "Starts the proxy server",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := startProxy(root.opts)
			if err != nil {
				return wrapError(err, fmt.Sprintln("Proxy server return an error"))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&root.opts.bindAddress, "bind-address", "b", ":8080", "Proxy listen address")

	root.cmd = cmd
	return root
}

func startProxy(options startOpts) (*context.Context, error) {
	ctx := context.New()
	ctx.BindAddress = options.bindAddress

	p := proxy.New(ctx)
	return ctx, p.Start()
}
