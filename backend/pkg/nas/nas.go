package nas

import (
	"bytes"
	"fmt"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/util"
	"os/exec"
	"strings"
)

// ZPool represents a ZFS zpool with relevant properties.
type ZPool struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	Allocated  string `json:"allocated"`
	Free       string `json:"free"`
	Fragmented string `json:"fragmented"`
	Health     string `json:"health"`
}

// ZFSDataset represents a ZFS volume with its name and quota.
type ZFSDataset struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Quota        string `json:"quota"`
	Used         string `json:"used"`
	Available    string `json:"available"`
	ShareEnabled bool   `json:"shareEnabled"`
}

// Snapshot represents the detailed information of a ZFS snapshot.
type Snapshot struct {
	Name       string `json:"name"`
	Used       string `json:"used"`
	Referenced string `json:"referenced"`
	CreatedAt  string `json:"createdAt"`
}

// ListZPools lists all zpools on the system.
func ListZPools() ([]ZPool, error) {
	cmd := exec.Command("zpool", "list", "-H", "-o", "name,size,alloc,free,frag,health")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var zpools []ZPool
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		zpools = append(zpools, ZPool{
			Name:       fields[0],
			Size:       fields[1],
			Allocated:  fields[2],
			Free:       fields[3],
			Fragmented: fields[4],
			Health:     fields[5],
		})
	}
	return zpools, nil
}

// ListZFSDatasets lists all ZFS volumes on the system.
func ListZFSDatasets() ([]ZFSDataset, error) {
	cmd := exec.Command("zfs", "list", "-H", "-o", "name,quota,used,avail", "-t", "filesystem")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var datasets []ZFSDataset
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		datasets = append(datasets, ZFSDataset{
			ID:        util.Base64Encode(fields[0]),
			Name:      fields[0],
			Quota:     fields[1],
			Used:      fields[2],
			Available: fields[3],
		})
	}
	return datasets, nil
}

// ListZVOLs lists all ZFS ZVOLs.
func ListZVOLs() ([]string, error) {
	cmd := exec.Command("zfs", "list", "-H", "-o", "name", "-t", "volume")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// CreateZFSVolume creates a ZFS volume with a specified quota.
func CreateZFSVolume(name, quota string) error {
	cmd := exec.Command("zfs", "create", "-o", fmt.Sprintf("quota=%s", quota), name)
	return cmd.Run()
}

// UpdateQuota updates the quota for an existing ZFS volume.
func UpdateQuota(volumeName, quota string) error {
	cmd := exec.Command("zfs", "set", fmt.Sprintf("quota=%s", quota), volumeName)
	return cmd.Run()
}

// CreateNFSShare creates an NFS share for a given ZFS volume.
func CreateNFSShare(zfsDatasetName string, rwIPs []string, roIPs []string) error {
	var rwAccess, roAccess string

	if len(rwIPs) > 0 {
		rwAccess = fmt.Sprintf(",rw=%s", strings.Join(rwIPs, ":"))
	}
	if len(roIPs) > 0 {
		roAccess = fmt.Sprintf(",ro=%s", strings.Join(roIPs, ":"))
	}

	shareNfs := fmt.Sprintf("insecure%s%s", rwAccess, roAccess)

	log.Logger.Infow("creating nfs share", "permission", shareNfs)

	cmd := exec.Command("zfs", "set", fmt.Sprintf("sharenfs=%s", shareNfs), zfsDatasetName)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return setPathPermissions(fmt.Sprintf("/%s", zfsDatasetName))
}

// setPathPermissions sets ownership to nobody:nogroup and permissions to 777 on the specified ZFS path.
func setPathPermissions(zfsPath string) error {
	// Change ownership to nobody:nogroup
	chownCmd := exec.Command("sudo", "chown", "-R", "nobody:nogroup", zfsPath)
	if err := chownCmd.Run(); err != nil {
		return fmt.Errorf("failed to set path ownership: %v", err)
	}

	// Change permissions to 777
	/*chmodCmd := exec.Command("sudo", "chmod", "-R", "777", zfsPath)
	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to set path permissions: %v", err)
	}*/

	return nil
}

// ConfigureNFSAccess sets NFS access for specified IPs with given permissions.
func ConfigureNFSAccess(volumeName string, rwIPs []string, roIPs []string) error {
	rwAccess := fmt.Sprintf("rw=%s", strings.Join(rwIPs, ":"))
	roAccess := fmt.Sprintf("ro=%s", strings.Join(roIPs, ":"))
	shareNfs := fmt.Sprintf("%s,%s,insecure", rwAccess, roAccess)

	cmd := exec.Command("zfs", "set", fmt.Sprintf("sharenfs=%s", shareNfs), volumeName)
	return cmd.Run()
}

// RemoveNFSAccess revokes an IP's access to a specified NFS share.
func RemoveNFSAccess(volumeName, ip string) error {
	cmd := exec.Command("zfs", "set", "sharenfs=off", volumeName)
	return cmd.Run()
}

// DeleteNFSShare disables NFS sharing on a ZFS volume.
func DeleteNFSShare(volumeName string) error {
	cmd := exec.Command("zfs", "set", "sharenfs=off", volumeName)
	return cmd.Run()
}

// DeleteZFSVolume deletes a specified ZFS volume.
func DeleteZFSVolume(volumeName string) error {
	cmd := exec.Command("zfs", "destroy", volumeName)
	return cmd.Run()
}

// ListSnapshots lists all snapshots for a given ZFS dataset with detailed information.
func ListSnapshots(dataset string) ([]Snapshot, error) {
	// Execute the zfs command to list snapshots with additional fields
	cmd := exec.Command("zfs", "list", "-t", "snapshot", "-o", "name,used,referenced,creation", "-H", "-d", "1", dataset)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}

	// Parse the output into Snapshot objects
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var snapshots []Snapshot
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue // Skip malformed lines
		}
		snapshots = append(snapshots, Snapshot{
			Name:       fields[0],
			Used:       fields[1],
			Referenced: fields[2],
			CreatedAt:  strings.Join(fields[3:], " "), // Creation time can have spaces
		})
	}

	return snapshots, nil
}

// CreateSnapshot creates a snapshot for a given ZFS dataset.
func CreateSnapshot(dataset, snapshotName string) error {
	snapshot := fmt.Sprintf("%s@%s", dataset, snapshotName)
	// Execute the zfs command to create the snapshot
	cmd := exec.Command("zfs", "snapshot", snapshot)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}
	return nil
}

// RestoreFromSnapshot rolls back a dataset to a given snapshot.
func RestoreFromSnapshot(snapshotName string) error {
	// Execute the zfs command to rollback the dataset to the snapshot
	cmd := exec.Command("zfs", "rollback", "-r", snapshotName)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
	if err != nil {
		return fmt.Errorf("failed to restore from snapshot '%s': %s (%w)", snapshotName, string(output), err)
	}
	return nil
}

// DeleteSnapshot deletes a specific ZFS snapshot.
func DeleteSnapshot(snapshotName string) error {
	// Execute the zfs destroy command
	cmd := exec.Command("zfs", "destroy", snapshotName)
	output, err := cmd.CombinedOutput() // Capture both stdout and stderr
	if err != nil {
		return fmt.Errorf("failed to delete snapshot '%s': %s (%w)", snapshotName, string(output), err)
	}
	return nil
}
