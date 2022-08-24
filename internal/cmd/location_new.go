package cmd

import (
	"fmt"

	"github.com/0ghny/gitconfig/pkg/locations"
	"github.com/0ghny/go-libx/pkg/iox"
	"github.com/spf13/cobra"
)

const (
	varLocation      = "location"
	varShortLocation = "l"
)

var (
	location string
)

func locationNewCmd() *cobra.Command {
	locationNewCmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new gitconfig for specified location",
		Long:  `Creates a new gitconfig file for specified location`,
		RunE: func(cmd *cobra.Command, args []string) error {
			locationMgr := locations.NewLocationManager("", nil)
			// Save location
			err := locationMgr.SaveLocation(inputLocationKey, location)
			if err != nil {
				// Error saving location
				return err
			}
			// Read new location
			location, err := locationMgr.FindLocationByKey(inputLocationKey)
			if err != nil {
				// error reading just saved location
				return err
			}
			fmt.Println(
				fmt.Sprintf("Location %s saved successfully, configuration file was created at %s", location.Key, location.ConfigFile))
			return nil
		},
	}
	addlocationNewCmdFlags(locationNewCmd)
	return locationNewCmd
}

func addlocationNewCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVarP(&location, varLocation, varShortLocation, iox.Getwd(), "Location for new gitconfig, default to current directory")
}
