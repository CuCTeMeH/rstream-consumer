package cmd

import "github.com/spf13/cobra"

func NewRootCMD() *cobra.Command {
	CMDInstance := &cobra.Command{
		Use:           "redis-stream-consumer",
		Short:         "redis stream consumer CLI tool",
		Long:          `redis stream consumer CLI tool`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	CMDInstance.AddCommand(NewConsumerCMD())

	return CMDInstance
}
