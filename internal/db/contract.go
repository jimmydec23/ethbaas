package db

import "fmt"

type Contract struct {
	Name    string `xorm:"name not null pk"`
	Proj    string `xorm:"proj"`
	Created int64  `xorm:"created not null"`
	Address string `xorm:"address"`
	ABI     string `xom:"path not null"`
	BIN     string `xom:"path not null"`
}

func (c *Client) AddContract(contract *Contract) error {
	_, err := c.engine.InsertOne(contract)
	return err
}

func (c *Client) ListContract() (int64, []Contract, error) {
	list := []Contract{}
	total, err := c.engine.FindAndCount(&list)
	return total, list, err
}

func (c *Client) GetContract(name string) (*Contract, error) {
	contract := &Contract{Name: name}
	has, err := c.engine.Get(contract)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("Contract %s not found", name)
	}
	return contract, nil
}
