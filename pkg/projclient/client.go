package projclient

import (
	"ethbaas/internal/db"
	"ethbaas/internal/k8s"
	"ethbaas/internal/model"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	db     *db.Client
	parser *k8s.Parser
}

func NewClient(db *db.Client) *Client {
	c := &Client{
		db:     db,
		parser: k8s.NewParser(),
	}
	return c
}

// init project, generate project yaml templates
func (c *Client) Init(p *model.Project) error {
	if c.db.IsProjectExist(&db.Project{Name: p.Name}) {
		return fmt.Errorf("Project %s already exist.", p.Name)
	}

	if err := c.parser.Parse(p); err != nil {
		return err
	}
	dbProj := &db.Project{
		Name:      p.Name,
		NodeCount: p.NodeCount,
		Created:   time.Now().Unix(),
		NodePort:  p.Port2Str(),
	}

	if err := c.db.AddProject(dbProj); err != nil {
		return err
	}

	return nil
}

// list projects
func (c *Client) List() (int64, []db.Project, error) {
	return c.db.ListProject()
}

// start project
func (c *Client) Start(projName string) error {
	dbproj, err := c.db.GetProject(projName)
	if err != nil {
		return err
	}

	proj := &model.Project{
		Name:      dbproj.Name,
		NodeCount: dbproj.NodeCount,
	}

	info, err := ioutil.ReadDir(proj.Home())
	if err != nil {
		log.Fatal(err)
	}

	yamlFiles := []string{}
	for _, file := range info {
		if !file.IsDir() {
			full := filepath.Join(proj.Home(), file.Name())
			yamlFiles = append(yamlFiles, full)
		}
	}
	for _, p := range yamlFiles {
		err := k8s.Apply(p)
		if err != nil {
			return err
		}
	}

	dbproj.Running = true
	if err := c.db.UpdateProject(dbproj); err != nil {
		return err
	}

	return nil
}

// stop project
func (c *Client) Stop(projName string) error {
	dbproj, err := c.db.GetProject(projName)
	if err != nil {
		return err
	}

	proj := &model.Project{
		Name:      dbproj.Name,
		NodeCount: dbproj.NodeCount,
	}

	info, err := ioutil.ReadDir(proj.Home())
	if err != nil {
		return err
	}

	yamlFiles := []string{}
	for _, file := range info {
		if !file.IsDir() {
			full := filepath.Join(proj.Home(), file.Name())
			yamlFiles = append(yamlFiles, full)
		}
	}
	for _, p := range yamlFiles {
		err := k8s.Delete(p)
		if err != nil {
			return nil
		}
	}

	dbproj.Running = false
	if err := c.db.UpdateProject(dbproj); err != nil {
		return err
	}
	return nil
}

// delete project
func (c *Client) Delete(projName string) error {
	dbproj, err := c.db.GetProject(projName)
	if err != nil {
		return err
	}

	proj := &model.Project{
		Name:      dbproj.Name,
		NodeCount: dbproj.NodeCount,
	}
	if err := os.RemoveAll(proj.Home()); err != nil {
		return err
	}

	dbProj := &db.Project{
		Name: projName,
	}
	if err := c.db.DeleteProject(dbProj); err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(projName string) (*db.Project, error) {
	return c.db.GetProject(projName)
}

func (c *Client) GetInModel(projName string) (*model.Project, error) {
	dbProj, err := c.Get(projName)
	if err != nil {
		return nil, err
	}
	m := &model.Project{
		Name:          dbProj.Name,
		NodeCount:     dbProj.NodeCount,
		FirstNodePort: dbProj.Str2Port()[0],
	}
	return m, nil
}
