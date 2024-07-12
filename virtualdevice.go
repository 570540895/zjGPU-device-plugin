package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type VirtualDevice struct {
	HostPath    string  `json:"hostPath"`
	ContainerPath  string  `json:"containerPath"`
	Permission     string  `json:"permission"`
}

func (d *VirtualDevice) validate() error {
	numGlobHostPath := strings.Count(d.HostPath, "*")
	if numGlobHostPath > 1 {
		return fmt.Errorf("HostPath can container only one '*' character: %s", d.HostPath)
	}

	if numGlobVirtualPath == 1 {
		if !strings.HasSuffix(d.ContainerPath, "*") {
			return fmt.Errorf("ContainerPath should ends with '*' character when VirtualPath container '*': %s", d.ContainerPath)
		}
		return nil
	}

	if strings.Contains(d.ContainerPath, "*") {
		return fmt.Errorf("ContainerPath must not contain '*' when VirtualPath does not contain '*': %s", d.ContainerPath)
	}

	return nil
}

type ExpandedVirtualDevice struct {
	HostPath      string
	ContainerPath string
	Permission    string
}

func (d VirtualDevice) Expand() ([]*ExpandedVirtualDevice, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}

	matchedHostPath, err := filepath.Glob(d.HostPath)
	if err != nil {
		return nil, err
	}

	expanded := []*ExpandedVirtualDevice{}
	baseHostPath := strings.Split(d.HostPath, "*")[0]
	baseContainerPath := strings.Split(d.ContainerPath, "*")[0]
	for _, hp := range matchedHostPath {
		fInfo, _ := os.Stat(hp)
		if fInfo.IsDir() {
			continue
		}

		expanded = append(expanded, &ExpandedVirtualDevice{
			HostPath:      hp,
			ContainerPath: strings.Replace(hp, baseHostPath, baseContainerPath, 1),
			Permission:    d.Permission,
		})
	}
	return expanded, nil
}
