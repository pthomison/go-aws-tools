package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg/client"
)

const (
	moveVolumeNameFlag      = "volume-name"
	moveVolumeNameShortFlag = "v"

	moveVolumeIdFlag      = "volume-id"
	moveVolumeIdShortFlag = "i"

	moveNewAvailabilityZoneFlag      = "availability-zone"
	moveNewAvailabilityZoneShortFlag = "a"

	// moveVolumeTypeFlag      = "type"
	// moveVolumeTypeShortFlag = "t"
)

var moveVolumeCmd = &cobra.Command{
	Use:   "move-volume",
	Short: "Utility to move EBS volumes between Availablity Zones",
	Long:  ``,
	Run:   moveVolumeCobra,
	Args:  cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(moveVolumeCmd)
	moveVolumeCmd.PersistentFlags().StringP(moveVolumeNameFlag, moveVolumeNameShortFlag, "", "Name of volume to move (mutually exlusive from "+moveVolumeIdFlag+")")
	moveVolumeCmd.PersistentFlags().StringP(moveVolumeIdFlag, moveVolumeIdShortFlag, "", "ID of volume to move (mutually exlusive from "+moveVolumeNameFlag+")")
	moveVolumeCmd.PersistentFlags().StringP(moveNewAvailabilityZoneFlag, moveNewAvailabilityZoneShortFlag, "", "AZ to move the volume to")
}

func moveVolumeCobra(cmd *cobra.Command, args []string) {
	volumeNameF := cmd.Flags().Lookup(moveVolumeNameFlag)
	volumeIdF := cmd.Flags().Lookup(moveVolumeIdFlag)
	newAvailabilityZoneF := cmd.Flags().Lookup(moveNewAvailabilityZoneFlag)

	// flag checking
	mutualExclusiveFlag(cmd, volumeNameF, volumeIdF)
	requiredFlag(cmd, newAvailabilityZoneF)

	volumeId, err := cmd.Flags().GetString(moveVolumeIdFlag)
	commandError(err)

	newAZ, err := cmd.Flags().GetString(moveNewAvailabilityZoneFlag)
	commandError(err)

	// initialize client
	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
	commandError(err)

	// lookup volume
	oldVol, err := client.FindVolume(volumeId)
	commandError(err)
	fmt.Printf("Volume: %+v\n", vol)

	// create snapshot
	snapshot, err := client.SnapshotVolume(vol)
	commandError(err)
	fmt.Printf("Snapshot: %+v\n", snapshot)

	// create volume from snapshot
	newVol, err := client.CreateVolumeFromSnapshot(snapshot, vol, newAZ)
	commandError(err)
	fmt.Printf("newVol: %+v\n", newVol)

	// rename old volume
	err = client.NameResource(volumeId, *tagLookup(oldVol.Tags, "Name").Value + "-old")
	commandError(err)
}
