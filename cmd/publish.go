package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/alexhokl/helper/iohelper"
	"github.com/spf13/cobra"
)

type publishOptions struct {
	file    string
	project string
	topic   string
}

var (
	publishOpts publishOptions
)

var publishCmd = &cobra.Command{
	Use:     "publish",
	Short:   "Publish a message",
	PreRunE: validatePublishOptions,
	RunE:    runPublish,
}

func init() {
	rootCmd.AddCommand(publishCmd)

	flags := publishCmd.Flags()
	flags.StringVarP(&publishOpts.file, "file", "f", "", "File to publish")
	flags.StringVarP(&publishOpts.project, "project", "p", "", "Project ID")
	flags.StringVarP(&publishOpts.topic, "topic", "t", "", "Topic name")

	publishCmd.MarkFlagRequired("file")
	publishCmd.MarkFlagRequired("project")
	publishCmd.MarkFlagRequired("topic")
}

func validatePublishOptions(_ *cobra.Command, _ []string) error {
	if publishOpts.file != "" {
		if _, err := os.Stat(publishOpts.file); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file [%s] does not exist: %v", publishOpts.file, err)
		}
		return nil
	}

	if publishOpts.project == "" {
		return fmt.Errorf("project is required")
	}

	if publishOpts.topic == "" {
		return fmt.Errorf("topic is required")
	}

	return nil
}

func runPublish(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, publishOpts.project)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	topic := client.Topic(publishOpts.topic)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check topic: %w", err)
	}
	if !ok {
		return fmt.Errorf("topic does not exist")
	}
	defer topic.Stop()

	fileBytes, err := iohelper.ReadBytesFromFile(publishOpts.file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	publishResult := topic.Publish(ctx, &pubsub.Message{
		Data:        fileBytes,
	})
	if _, err := publishResult.Get(ctx); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Println("Message published")
	return nil
}
