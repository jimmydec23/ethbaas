package db

import (
	"fmt"
	"strconv"
	"strings"
)

type Project struct {
	Name      string `xorm:"name varchar(25) not null pk"`
	NodeCount int    `xorm:"nodeCount not null"`
	Running   bool   `xorm:"running"`
	Created   int64  `xorm:"created not null"`
	NodePort  string `xorm:"nodePort not null"`
}

func (p *Project) Str2Port() []int32 {
	ports := strings.Split(p.NodePort, ",")
	portInt32 := []int32{}
	for _, p := range ports {
		pi, _ := strconv.Atoi(p)
		portInt32 = append(portInt32, int32(pi))
	}
	return portInt32
}

func (c *Client) AddProject(p *Project) error {
	_, err := c.engine.InsertOne(p)
	return err
}

func (c *Client) DeleteProject(p *Project) error {
	_, err := c.engine.Delete(p)
	return err
}

func (c *Client) IsProjectExist(p *Project) bool {
	has, _ := c.engine.Get(p)
	return has
}

func (c *Client) UpdateProject(p *Project) error {
	_, err := c.engine.UseBool().Where("name = ?", p.Name).Update(p)
	return err
}

func (c *Client) GetProject(name string) (*Project, error) {
	p := &Project{Name: name}
	has, err := c.engine.Get(p)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("Project %s not found.", name)
	}
	return p, nil
}

func (c *Client) ListProject() (int64, []Project, error) {
	list := []Project{}
	total, err := c.engine.FindAndCount(&list)
	return total, list, err
}
