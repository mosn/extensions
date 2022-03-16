package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"sync"

	"mosn.io/extensions/go-plugin/pkg/transcoder/bumscd"
)

type config struct {
	lock     sync.RWMutex
	relation *bumscd.Relation
	md5      []byte
}

func NewConfig(info string) (*config, error) {
	relation := &bumscd.Relation{}
	if err := json.Unmarshal(bumscd.S2B(info), relation); err != nil {
		return nil, err
	}
	return &config{
		md5:      md5.New().Sum(bumscd.S2B(info)),
		relation: relation,
	}, nil
}

func (c *config) update(relation *bumscd.Relation, md5 []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.relation = relation
	c.md5 = md5
}

func (c *config) GetLatestRelation(info string) (*bumscd.Relation, error) {
	c.lock.RLock()
	relation := &bumscd.Relation{}
	rmd5 := md5.New().Sum(bumscd.S2B(info))
	if bytes.Compare(rmd5, c.md5) == 0 {
		return c.relation, nil
	}

	if err := json.Unmarshal(bumscd.S2B(info), relation); err != nil {
		return nil, err
	}
	c.md5 = rmd5
	c.relation = relation
	return relation, nil
}

func ParseRelation(cfg map[string]interface{}) (*bumscd.Relation, error) {
	rInfo, ok := cfg["relation"]
	if !ok {
		return nil, nil
	}
	info, ok := rInfo.(string)
	if !ok {
		return nil, nil
	}

	uinfo, ok := cfg["unique"]
	if !ok {
		return nil, nil
	}
	uid, ok := uinfo.(string)
	if !ok {
		return nil, nil
	}
	if config, ok := relations[uid]; ok {
		return config.relation, nil
	}

	config, err := NewConfig(info)
	if err != nil {
		return nil, err
	}
	relations[uid] = config
	return config.relation, nil
}
