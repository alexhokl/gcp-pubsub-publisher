package cmd

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	project      string
	subscription string
	acknowledge  bool
}

var (
	pullOpts pullOptions
)

var pullCmd = &cobra.Command{
	Use:     "pull",
	Short:   "Pull a message",
	PreRunE: validatePullOptions,
	RunE:    runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)

	flags := pullCmd.Flags()
	flags.StringVarP(&pullOpts.project, "project", "p", "", "Project ID")
	flags.StringVarP(&pullOpts.subscription, "subscription", "s", "", "Subscription name")
	flags.BoolVarP(&pullOpts.acknowledge, "acknowledge", "a", false, "Acknowledge the message")

	pullCmd.MarkFlagRequired("project")
	pullCmd.MarkFlagRequired("subscription")
}

func validatePullOptions(_ *cobra.Command, _ []string) error {
	if pullOpts.project == "" {
		return fmt.Errorf("project is required")
	}

	if pullOpts.subscription == "" {
		return fmt.Errorf("subscription is required")
	}

	return nil
}

func runPull(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, pullOpts.project)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	sub := client.Subscription(pullOpts.subscription)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check subscription: %w", err)
	}
	if !ok {
		return fmt.Errorf("subscription does not exist")
	}

	sub.ReceiveSettings.MaxOutstandingMessages = 1

	fmt.Printf("Listening for messages on %s\n", pullOpts.subscription)

	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Printf("message %s: %s\n", msg.ID, string(msg.Data))
		if pullOpts.acknowledge {
			msg.Ack()
		} else {
			msg.Nack()
		}

		fmt.Println("message pulled")
	})
	if err != nil {
		return fmt.Errorf("failed to receive message: %w", err)
	}

	return nil
}
