package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type VirtualDevice struct {
	VirtualPath    string  `json:"virtualPath"`
	ContainerPath  string  `json:"containerPath"`
	Permission     string  `json:"permission"`
}

func (d *VirtualDevice) validate() error {
	numGlobVirtualPath := strings.Count(d.VirtualPath, "*")
	if numGlobVirtualPath > 1 {
		return fmt.Errorf("VirtualPath can container only one '*' character: %s", d.VirtualPath)
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
	VirtualPath      string
	ContainerPath string
	Permission    string
}

func (d VirtualDevice) Expand() ([]*ExpandedVirtualDevice, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}

	matchedVirtualPath, err := filepath.Glob(d.VirtualPath)
	if err != nil {
		return nil, err
	}

	expanded := []*ExpandedVirtualDevice{}
	baseVirtualPath := strings.Split(d.VirtualPath, "*")[0]
	baseContainerPath := strings.Split(d.ContainerPath, "*")[0]
	for _, vp := range matchedVirtualPath {
		fInfo, _ := os.Stat(vp)
		if fInfo.IsDir() {
			continue
		}

		expanded = append(expanded, &ExpandedVirtualDevice{
			VirtualPath:      vp,
			ContainerPath: strings.Replace(hp, baseVirtualPath, baseContainerPath, 1),
			Permission:    d.Permission,
		})
	}
	return expanded, nil
}
