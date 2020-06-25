package client

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Client) FindVolume(id string) (*ec2.Volume, error) {
	svc := ec2.New(c.sess)

	vol, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(id)},
	})

	if err != nil {
		return nil, err
	}

	return vol.Volumes[0], nil
}

// func (c *Client) FindVolumeByName(volumeId string) (*ec2.Volume, error) {
// }

func (c *Client) SnapshotVolume(vol *ec2.Volume) (*ec2.Snapshot, error) {
	svc := ec2.New(c.sess)

	volId := vol.VolumeId
	volName := tagLookup(vol.Tags, "Name").Value

	tagSpec := &ec2.TagSpecification{
		ResourceType: aws.String("snapshot"),
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(*volId + "-snapshot"),
			},
			&ec2.Tag{
				Key:   aws.String("SourceVolume"),
				Value: aws.String(*volName),
			},
			&ec2.Tag{
				Key:   aws.String("Reason"),
				Value: aws.String("inter_az_copy"),
			},
		},
	}

	snapshot, err := svc.CreateSnapshot(&ec2.CreateSnapshotInput{
		VolumeId:          volId,
		TagSpecifications: []*ec2.TagSpecification{tagSpec},
	})

	snapshot.SetDescription("Moving Volumes Between AZs")

	if err != nil {
		return nil, err
	}

	snapshotId := snapshot.SnapshotId

	for {
		time.Sleep(15 * time.Second)

		status, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
			SnapshotIds: []*string{snapshotId},
		})
		if err != nil {
			return nil, err
		}

		fmt.Printf("%+v\n", status.Snapshots[0])

		if *status.Snapshots[0].State == "completed" {
			break
		}
	}

	return snapshot, nil
}

func (c *Client) CreateVolumeFromSnapshot(snapshot *ec2.Snapshot, oldVolume *ec2.Volume, newAz string) (*ec2.Volume, error) {
	svc := ec2.New(c.sess)

	volumeName := *tagLookup(oldVolume.Tags, "Name").Value
	volumeType := *oldVolume.VolumeType

	snapshotId := *snapshot.SnapshotId

	tagSpec := &ec2.TagSpecification{
		ResourceType: aws.String("volume"),
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(volumeName),
			},
		},
	}

	volInput := &ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String(newAz),
		SnapshotId:        aws.String(snapshotId),
		VolumeType:        aws.String(volumeType),
		TagSpecifications: []*ec2.TagSpecification{tagSpec},
	}

	if volumeType == "io1" {
		volInput.Iops = oldVolume.Iops
	}

	fmt.Printf("%+v\n", result)

	return svc.CreateVolume(volInput)
}
