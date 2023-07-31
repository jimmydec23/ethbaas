package db

import (
	"database/sql"
	"ethbaas/internal/config"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type Client struct {
	db     *sql.DB
	engine *xorm.Engine
}

func NewClient() (*Client, error) {
	store := config.C.GetString("dbstore")

	engine, err := xorm.NewEngine("sqlite3", store)
	if err != nil {
		return nil, err
	}
	if config.C.GetBool("showsql") {
		engine.SetLogLevel(log.LOG_INFO)
		engine.ShowSQL(true)
	} else {
		engine.SetLogLevel(log.LOG_ERR)
	}
	if err := engine.Ping(); err != nil {
		return nil, err
	}

	c := &Client{
		engine: engine,
	}
	return c, nil
}

// setup
func (c *Client) Setup() error {
	return c.tableInitial()
}

// init tables
func (c *Client) tableInitial() error {
	t := &Project{}
	exist, err := c.engine.IsTableExist(t)
	if err != nil {
		return err
	}
	if !exist {
		c.engine.CreateTables(t)
	}

	t2 := &Contract{}
	exist, err = c.engine.IsTableExist(t2)
	if err != nil {
		return err
	}
	if !exist {
		c.engine.CreateTables(t2)
	}

	return nil
}

// close engine
func (c *Client) Close() {
	c.engine.Close()
}
