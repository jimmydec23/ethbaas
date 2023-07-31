package model

import (
	"ethbaas/internal/config"
	"fmt"
	"path/filepath"
	"strings"
)

type Project struct {
	Name          string
	NodeCount     int
	FirstNodePort int32
}

func (p *Project) Home() string {
	home := filepath.Join(
		config.C.GetString("homedir"),
		p.Name,
	)
	return home
}

func (p *Project) NS() string {
	return fmt.Sprintf("ethbaas-%s", p.Name)
}

func (p *Project) NsFile() string {
	return filepath.Join(p.Home(), "1.ns.yaml")
}

func (p *Project) CmFile() string {
	return filepath.Join(p.Home(), "cm.yaml")
}

func (p *Project) PvFile(i int) string {
	return filepath.Join(p.Home(), fmt.Sprintf("pv_%d.yaml", i))
}

func (p *Project) PvcFile(i int) string {
	return filepath.Join(p.Home(), fmt.Sprintf("pvc_%d.yaml", i))
}

func (p *Project) SvcFile(i int) string {
	return filepath.Join(p.Home(), fmt.Sprintf("svc_%d.yaml", i))
}

func (p *Project) DeployFile(i int) string {
	return filepath.Join(p.Home(), fmt.Sprintf("deploy_%d.yaml", i))
}

func (p *Project) Port2Str() string {
	ports := []string{}
	for i, start := 0, p.FirstNodePort; i < p.NodeCount; i, start = i+1, start+1 {
		ports = append(ports, fmt.Sprintf("%d", start))
	}
	return strings.Join(ports, ",")
}
