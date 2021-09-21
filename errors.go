// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import (
	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

type StanzaErr struct {
	Stanza xc.Stanza
	Err    xc.Err
}

func (se StanzaErr) ToElem(elem *stravaganza.Element) {
	var err stravaganza.Element
	se.Err.ToElem(&err)
	*elem = se.Stanza.ToElemBuilder().WithChild(err).Build()
}

func (se *StanzaErr) FromElem(elem stravaganza.Element, name string) error {
	if err := se.Stanza.FromElem(elem, name); err != nil {
		return err
	}
	if err := elem.Child("error"); err != nil {
		se.Err.FromElem(elem.Child("error"))
	}
	return nil
}

type IqErrHandler interface {
	HandleIqError(StanzaErr, string) error
}

type StanzaErrHandler struct {
	Name             string
	lastMatchedError StanzaErr
	handler          IqErrHandler
}

func (seh *StanzaErrHandler) Match(elem stravaganza.Element) bool {
	if err := seh.lastMatchedError.FromElem(elem, seh.Name); err != nil {
		return false
	}
	return true
}

func (seh StanzaErrHandler) Handle(_ stravaganza.Element, part xc.Part) error {
	return seh.handler.HandleIqError(seh.lastMatchedError, seh.Name)
}
