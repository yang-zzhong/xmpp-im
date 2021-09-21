// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import (
	"errors"

	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

var (
	ErrNotIqErr = errors.New("not iq error")
)

type IqQuery struct {
	Stanza   xc.Stanza
	Xmlns    string
	Children []stravaganza.Element
	Ver      string
}

func (iqq IqQuery) ToElem(elem *stravaganza.Element) {
	iqq.Stanza.ToElemBuilder()
	query := stravaganza.NewBuilder("query").WithAttribute("xmlns", iqq.Xmlns)
	if len(iqq.Children) > 0 {
		query.WithChildren(iqq.Children...)
	}
	if iqq.Ver != "" {
		query.WithAttribute("ver", iqq.Ver)
	}
	*elem = iqq.Stanza.ToElemBuilder().WithChild(query.Build()).Build()
}

func (iqq *IqQuery) FromElem(elem stravaganza.Element) error {
	if err := iqq.Stanza.FromElem(elem, xc.NameIQ); err != nil {
		return err
	}
	query := elem.Child("query")
	if query == nil {
		return errors.New("not iq query element")
	}
	iqq.Xmlns = query.Attribute("xmlns")
	iqq.Ver = query.Attribute("ver")
	iqq.Children = query.AllChildren()
	return nil
}
