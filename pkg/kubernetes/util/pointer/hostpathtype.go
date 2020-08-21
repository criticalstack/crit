package pointer

import (
	"os"

	corev1 "k8s.io/api/core/v1"
)

func HostPathTypePtr(h corev1.HostPathType) *corev1.HostPathType {
	return &h
}

// DetectHostPathType attempts to determine the type of HostPathType for a
// given file. It currently only distinguishes between files and directories.
func DetectHostPathType(path string) (*corev1.HostPathType, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return HostPathTypePtr(corev1.HostPathDirectoryOrCreate), nil
	}
	return HostPathTypePtr(corev1.HostPathFileOrCreate), nil
}
