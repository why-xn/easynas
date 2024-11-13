package nas

import (
	"fmt"
	"os/exec"
	"strings"
)

// Zpool represents a ZFS zpool with relevant properties.
type Zpool struct {
	Name       string
	Size       string
	Allocated  string
	Free       string
	Fragmented string
	Health     string
}

// ZFSVolume represents a ZFS volume with its name and quota.
type ZFSVolume struct {
	Name      string
	Quota     string
	Used      string
	Available string
}

// ListZpools lists all zpools on the system.
func ListZpools() ([]Zpool, error) {
	cmd := exec.Command("zpool", "list", "-H", "-o", "name,size,alloc,free,frag,health")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var zpools []Zpool
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		zpools = append(zpools, Zpool{
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

// ListZFSVolumes lists all ZFS volumes on the system.
func ListZFSVolumes() ([]ZFSVolume, error) {
	cmd := exec.Command("zfs", "list", "-H", "-o", "name,quota,used,avail", "-t", "filesystem")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var volumes []ZFSVolume
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		volumes = append(volumes, ZFSVolume{
			Name:      fields[0],
			Quota:     fields[1],
			Used:      fields[2],
			Available: fields[3],
		})
	}
	return volumes, nil
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
func CreateNFSShare(poolName, volumeName string) error {
	cmd := exec.Command("zfs", "set", "sharenfs=on", fmt.Sprintf("%s/%s", poolName, volumeName))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return setPathPermissions(fmt.Sprintf("%s/%s", poolName, volumeName))
}

// setPathPermissions sets ownership to nobody:nogroup and permissions to 777 on the specified ZFS path.
func setPathPermissions(zfsPath string) error {
	// Change ownership to nobody:nogroup
	chownCmd := exec.Command("sudo", "chown", "-R", "nobody:nogroup", zfsPath)
	if err := chownCmd.Run(); err != nil {
		return fmt.Errorf("failed to set path ownership: %v", err)
	}

	// Change permissions to 777
	chmodCmd := exec.Command("sudo", "chmod", "-R", "777", zfsPath)
	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to set path permissions: %v", err)
	}

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