package cmd

import (
	"fmt"
	"os"

	"github.com/0ghny/gitconfig/pkg/locations"
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	varLocationKey      = "key"
	varShortLocationKey = "k"
)

var (
	inputLocationKey string
)

func locationCmd() *cobra.Command {
	locationCmd := &cobra.Command{
		Use:   "location",
		Short: "Manage a gitconfig location",
		Long:  `Manage a gitconfig location`,
		RunE: func(cmd *cobra.Command, args []string) error {
			locationMgr := locations.NewLocationManager("", nil)
			location, err := locationMgr.FindLocationByKey(inputLocationKey)
			if err != nil {
				log.Error(err)
				return err
			}

			// location not found
			if location == nil {
				log.Debug("location %s not found", inputLocationKey)
				return fmt.Errorf("location %s not found", inputLocationKey)
			}

			// found! printing
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"#", "Key", "Location", "GitConfig"})
			t.AppendRow([]interface{}{"#", location.Key, location.Path, location.ConfigFile})
			t.AppendSeparator()
			t.Render()
			return nil
		},
	}

	// Configure
	addLocationCmdFlags(locationCmd)
	// Subcommands
	locationCmd.AddCommand(locationNewCmd())
	return locationCmd
}

func addLocationCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().
		StringVarP(&inputLocationKey, varLocationKey, varShortLocationKey, "", "Location key, it will be used to identify the location in further operations")
	cmd.MarkPersistentFlagRequired(varLocationKey)
}
