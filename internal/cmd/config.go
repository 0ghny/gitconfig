package cmd

import (
	"errors"
	"fmt"

	"github.com/0ghny/gitconfigs/internal/git"
	"github.com/0ghny/gitconfigs/pkg/locations"
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
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
			// Get mandatory argument (config key)
			configKey := args[0]

			// If not Key is provided, then, run the command just as regular git config
			// in this case, it's just a wrapper
			gitConfigfile := ""
			if inputLocationKey != "" {
				// Retrieve location info from inputLocation (mandatory parameter)
				locationMgr := locations.NewLocationManager("", nil)
				// Read specified location
				location, err := locationMgr.FindLocationByKey(inputLocationKey)
				if err != nil {
					// error reading just saved location
					return err
				}
				if location == nil {
					return fmt.Errorf("location with key %s not found", inputLocationKey)
				}
				gitConfigfile = location.ConfigFile
			}

			// Two options, if 2 arguments, or only one provided
			if len(args) == 1 {
				// GET
				out, err := git.GitConfigGet(configKey, gitConfigfile)
				if err != nil {
					//if error, most probably is that config key doens't exists
					return fmt.Errorf("config key `%s` doesn't exists in file `%s`", configKey, gitConfigfile)
				}
				fmt.Println(out)
			} else if len(args) == 2 {
				// SET
				configValue := args[1]
				err := git.GitConfigSet(configKey, configValue, gitConfigfile)
				if err != nil {
					return err
				}
				fmt.Printf("config set successfully in config file %s\n", gitConfigfile)
			} else {
				return errors.New("minium argument is 1 (config key) and maximum 2 (config key + value to set)")
			}

			return nil
		},
	}
	// Configure
	addConfigCmdFlags(configCmd)
	return configCmd
}

func addConfigCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVarP(&inputLocationKey, varLocationKey, varShortLocationKey, "", "Location key, it will be used to identify the location in further operations")
}
