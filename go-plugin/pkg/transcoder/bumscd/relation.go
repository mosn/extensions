package bumscd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

type Relation struct {
	SysHead   []*BumsAndCdIterm          `json:"sys_head"`
	AppHead   []*BumsAndCdIterm          `json:"app_head"`
	LocalHead []*BumsAndCdIterm          `json:"local_head"`
	Body      map[string]*BumsAndCdIterm `json:"body"`
	IsResp    bool                       `json:"-"`
}

func NewRelation(isResp bool) *Relation {
	return &Relation{
		IsResp: isResp,
	}
}

func (rl *Relation) ParseFile(headPath, bodyPath string) error {
	root := etree.NewDocument()
	if err := root.ReadFromFile(headPath); err != nil {
		return err
	}
	if err := rl.ParseHead(root); err != nil {
		return err
	}

	broot := etree.NewDocument()
	if err := broot.ReadFromFile(bodyPath); err != nil {
		return err
	}
	if err := rl.ParseBody(broot); err != nil {
		return err
	}
	return nil
}

func (rl *Relation) ParseString(head, body string) error {
	root := etree.NewDocument()
	if err := root.ReadFromString(head); err != nil {
		return err
	}
	if err := rl.ParseHead(root); err != nil {
		return err
	}

	broot := etree.NewDocument()
	if err := broot.ReadFromString(body); err != nil {
		return err
	}
	if err := rl.ParseBody(broot); err != nil {
		return err
	}
	return nil
}

func (rl *Relation) ParseBody(root *etree.Document) (err error) {
	body := root.SelectElement("root")
	if body == nil {
		return fmt.Errorf("body is empty")
	}
	if rl.Body, err = rl.VisitBodyElements(body); err != nil {
		return err
	}
	return nil
}

func (rl *Relation) ParseHead(root *etree.Document) (err error) {
	head := root.SelectElement("root")
	if head == nil {
		return fmt.Errorf("head is empty")
	}

	// sys
	element := head.SelectElement("SYS_HEAD")
	if element == nil {
		return fmt.Errorf("the %v of SYS_HEAD is empty", head)
	}
	if rl.SysHead, err = rl.VisitHeadElements(element); err != nil {
		return err
	}

	// app
	element = head.SelectElement("APP_HEAD")
	if element == nil {
		return fmt.Errorf("the %v APP_HEAD is empty", head)
	}
	if rl.AppHead, err = rl.VisitHeadElements(element); err != nil {
		return err
	}

	// local
	element = head.SelectElement("LOCAL_HEAD")
	if element == nil {
		return fmt.Errorf("the %v of LOCAL_HEAD is empty", head)
	}
	if rl.LocalHead, err = rl.VisitHeadElements(element); err != nil {
		return err
	}
	return nil
}

func (rl *Relation) VisitHeadElements(root *etree.Element) ([]*BumsAndCdIterm, error) {
	elements := root.ChildElements()
	rmap := make([]*BumsAndCdIterm, 0, len(elements))
	for _, e := range elements {
		attr := e.SelectAttr("type")
		if attr != nil && strings.Compare(attr.Value, FieldList) == 0 {
			childMap, err := rl.VisitHeadElements(e)
			if err != nil {
				return nil, err
			}
			iterm, err := NewBumsAndCdList(e, childMap, rl.IsResp)
			if err != nil {
				return nil, err
			}
			rmap = append(rmap, iterm)
			continue
		}
		iterm, err := NewBumsAndCdIterm(e, rl.IsResp)
		if err != nil {
			return nil, err
		}
		rmap = append(rmap, iterm)
	}
	return rmap, nil
}

func (rl *Relation) VisitBodyElements(root *etree.Element) (map[string]*BumsAndCdIterm, error) {
	elements := root.ChildElements()
	rmap := make(map[string]*BumsAndCdIterm)
	for _, e := range elements {
		attr := e.SelectAttr("type")
		if attr != nil && strings.Compare(attr.Value, FieldList) == 0 {
			childMap, err := rl.VisitBodyElements(e)
			if err != nil {
				return nil, err
			}
			iterm, err := NewBumsAndCdBodyList(e, childMap, rl.IsResp)
			if err != nil {
				return nil, err
			}
			if rl.IsResp {
				rmap[iterm.CdKey] = iterm
			} else {
				rmap[iterm.BumsKey] = iterm
			}
			continue
		}
		iterm, err := NewBumsAndCdBodyIterm(e, rl.IsResp)
		if err != nil {
			return nil, err
		}
		if rl.IsResp {
			rmap[iterm.CdKey] = iterm
		} else {
			rmap[iterm.BumsKey] = iterm
		}
	}
	return rmap, nil
}

func (rl *Relation) String() string {
	data, _ := json.Marshal(rl)
	return B2S(data)
}
