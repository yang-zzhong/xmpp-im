// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package xmppim

import (
	"github.com/google/uuid"
	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

type sub string

const (
	nsRoster = "jabber:iq:roster"

	rosterInit = iota
	rosterError
	rosterQueryResult
	rosterPush
)

type RosterQueryResultHandler interface {
	HandleRosterQueryResult([]Roster)
}

type RosterPushHandler interface {
	HandleRosterPush([]Roster)
}

type RosterHandler interface {
	RosterQueryResultHandler
	RosterPushHandler
}
type RosterClient struct {
	id              string
	jid             string
	rosterHandler   RosterHandler
	lastIqQuery     IqQuery
	lastMatchedType int
	*xc.IDAble
}

func NewRosterClient(jid string, handler RosterHandler) *RosterClient {
	return &RosterClient{id: uuid.New().String(), jid: jid, rosterHandler: handler, IDAble: xc.NewIDAble()}
}

func (rm *RosterClient) Query(sender xc.Sender) error {
	iqq := IqQuery{
		Stanza: xc.Stanza{
			Name: xc.NameIQ,
			ID:   rm.id,
			Type: xc.StanzaGet,
			From: rm.jid}, Xmlns: nsRoster}
	var elem stravaganza.Element
	iqq.ToElem(&elem)
	return sender.SendElement(elem)
}

func (rm *RosterClient) Set(sender xc.Sender, roster Roster) error {
	elems := []stravaganza.Element{roster.ToElem()}
	var elem stravaganza.Element
	iqq := IqQuery{
		Stanza: xc.Stanza{
			Name: xc.NameIQ,
			ID:   rm.id,
			Type: xc.StanzaSet,
			From: rm.jid}, Xmlns: nsRoster, Children: elems}
	iqq.ToElem(&elem)
	return sender.SendElement(elem)
}

func (rm *RosterClient) Match(elem stravaganza.Element) bool {
	if elem.Child("error") != nil {
		rm.lastMatchedType = rosterError
		return true
	}
	if err := rm.lastIqQuery.FromElem(elem); err != nil {
		return false
	}
	switch rm.lastIqQuery.Stanza.Type {
	case xc.StanzaSet:
		rm.lastMatchedType = rosterPush
		return true
	case xc.StanzaResult:
		rm.lastMatchedType = rosterQueryResult
		return true
	}
	return false
}

func (rm *RosterClient) Handle(elem stravaganza.Element, part xc.Part) error {
	switch rm.lastMatchedType {
	case rosterQueryResult:
		rm.handleQueryResult(elem.Child("query").AllChildren())
	case rosterPush:
		rm.handlePush(elem.Child("query").AllChildren())
	}
	rm.lastMatchedType = rosterInit
	return nil
}

func (rm *RosterClient) handleQueryResult(elems []stravaganza.Element) {
	rosters := []Roster{}
	for _, elem := range elems {
		roster := Roster{}
		roster.FromElem(elem)
		rosters = append(rosters, roster)
	}
	rm.rosterHandler.HandleRosterQueryResult(rosters)
}

func (rm *RosterClient) handlePush(elems []stravaganza.Element) {
	rosters := []Roster{}
	for _, elem := range elems {
		roster := Roster{}
		roster.FromElem(elem)
		rosters = append(rosters, roster)
	}
	rm.rosterHandler.HandleRosterPush(rosters)
}
