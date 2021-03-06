package main

import (
	"os"

	"github.com/coreos/torus/block"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var volumeCommand = &cobra.Command{
	Use:   "volume",
	Short: "manage volumes in the cluster",
	Run:   volumeAction,
}

var volumeDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "delete a volume in the cluster",
	Run:   volumeDeleteAction,
}

var volumeListCommand = &cobra.Command{
	Use:   "list",
	Short: "list volumes in the cluster",
	Run:   volumeListAction,
}

func init() {
	volumeCommand.AddCommand(volumeDeleteCommand)
	volumeCommand.AddCommand(volumeListCommand)
	volumeListCommand.Flags().BoolVarP(&outputAsCSV, "csv", "", false, "output as csv instead")
}

func volumeAction(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

func volumeListAction(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
		os.Exit(1)
	}
	mds := mustConnectToMDS()
	vols, _, err := mds.GetVolumes()
	if err != nil {
		die("error listing volumes: %v\n", err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	if outputAsCSV {
		table.SetBorder(false)
		table.SetColumnSeparator(",")
	} else {
		table.SetHeader([]string{"Volume Name", "Size"})
	}
	for _, x := range vols {
		table.Append([]string{
			x.Name,
			humanize.IBytes(x.MaxBytes),
		})
	}
	table.Render()
}

func volumeDeleteAction(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(1)
	}
	name := args[0]
	mds := mustConnectToMDS()
	vol, err := mds.GetVolume(name)
	if err != nil {
		die("cannot get volume %s (perhaps it doesn't exist): %v", name, err)
	}
	switch vol.Type {
	case "block":
		err = block.DeleteBlockVolume(mds, name)
	default:
		die("unknown volume type %s", vol.Type)
	}
	if err != nil {
		die("cannot delete volume: %v", err)
	}
}
