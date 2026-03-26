package cli

import (
	"errors"
	"fmt"

	"github.com/0ghny/gitconfig/internal/adapter/git"
	"github.com/spf13/cobra"
)

// configCmd returns the "config" subcommand. It wraps `git config` to get or
// set a key in the gitconfig file associated with an optional --key location.
func configCmd(factory LocationServiceFactory) *cobra.Command {
	var key string

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Executes `git config [config-key] [config value] on an specified location",
		Long:  "Executes `git config [config-key] [config value] on an specified location",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) <= 0 || len(args) > 2 {
				return errors.New("minium argument is 1 (config key) and maximum 2 (config key + value to set)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configKey := args[0]

			// If no --key is provided, run as a regular git config wrapper.
			gitConfigfile := ""
			if key != "" {
				gitconfigPath, _ := cmd.Root().PersistentFlags().GetString(varGitconfigPath)
				svc := factory(gitconfigPath)
				loc, err := svc.FindLocationByKey(key)
				if err != nil {
					return err
				}
				if loc == nil {
					return fmt.Errorf("location with key %s not found", key)
				}
				gitConfigfile = loc.ConfigFile
			}

			if len(args) == 1 {
				out, err := git.GitConfigGet(configKey, gitConfigfile)
				if err != nil {
					return fmt.Errorf("config key `%s` doesn't exists in file `%s`", configKey, gitConfigfile)
				}
				fmt.Println(out)
			} else if len(args) == 2 {
				configValue := args[1]
				if err := git.GitConfigSet(configKey, configValue, gitConfigfile); err != nil {
					return err
				}
				fmt.Printf("config set successfully in config file %s\n", gitConfigfile)
			} else {
				return errors.New("minium argument is 1 (config key) and maximum 2 (config key + value to set)")
			}

			return nil
		},
	}
	configCmd.PersistentFlags().StringVarP(&key, varLocationKey, varShortLocationKey, "", "Location key, used to identify the location in further operations")
	return configCmd
}
