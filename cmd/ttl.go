package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	tikverror "github.com/tikv/client-go/v2/error"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Get the TTL of a key",
	RunE:  ttlRunE,
}

func ttlRunE(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer client.Close(cmd.Context())

	ttl, err := client.TTL(cmd.Context(), []byte(args[0]))
	if err != nil {
		if errors.Is(err, tikverror.ErrNotExist) {
			fmt.Fprintf(os.Stdout, "key [%s] not exist\n", args[0])
			return nil
		}

		return fmt.Errorf("get: %w", err)
	}

	if _, err := fmt.Fprintln(os.Stdout, ttl); err != nil {
		return err
	}
	return nil
}
