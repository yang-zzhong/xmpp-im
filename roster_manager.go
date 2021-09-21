// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

const (
	nsStanza    = "urn:ietf:params:xml:ns:xmpp-stanzas"
	rosterQuery = iota
	rosterSet

	RMForbidden     = "forbidden"
	RMBadRequest    = "bad request"
	RMNotAcceptable = "not acceptable"
)

var (
	ErrRMForbidden     = errors.New(RMForbidden)
	ErrRMBadRequest    = errors.New(RMBadRequest)
	ErrRMNotAcceptable = errors.New(RMNotAcceptable)
)

func RosterErrorElem(id, jid, typ string, err error) stravaganza.Element {
	tagName := strings.ReplaceAll(err.Error(), " ", "-")
	var elem stravaganza.Element
	StanzaErr{Stanza: xc.Stanza{
		Name: xc.NameIQ,
		ID:   id,
		From: jid,
		Type: xc.StanzaError,
	}, Err: xc.Err{
		Type:    typ,
		DescTag: tagName}}.ToElem(&elem)
	return elem
}

type IsAuthed interface {
	IsUserAuthed(jid string) bool
}

type RosterManageHandler interface {
	Rosters(rosters *[]Roster) error
	AddOrUpdate(roster Roster) error
	Delete(roster Roster) error
}

type RosterManager struct {
	lastMatchType int
	lastIqQuery   IqQuery
	handler       RosterManageHandler
	auth          IsAuthed
	*xc.IDAble
}

func NewRosterManager(handler RosterManageHandler, auth IsAuthed) *RosterManager {
	return &RosterManager{handler: handler, auth: auth, IDAble: xc.NewIDAble()}
}

func (rsm *RosterManager) Push(part xc.Part, rosters []Roster) error {
	to := part.Attr().JID.String()
	return part.Channel().SendElement(rsm.queryResult(rosters, uuid.New().String(), to, xc.StanzaSet))
}

func (rsm *RosterManager) Match(elem stravaganza.Element) bool {
	if err := rsm.lastIqQuery.FromElem(elem); err != nil {
		return false
	}
	iqq := &rsm.lastIqQuery
	switch iqq.Stanza.Type {
	case xc.StanzaGet:
		rsm.lastMatchType = rosterQuery
		return true
	case xc.StanzaSet:
		rsm.lastMatchType = rosterSet
		return true
	}
	return false
}

func (rsm *RosterManager) Handle(elem stravaganza.Element, part xc.Part) error {
	defer func() { rsm.lastMatchType = rosterInit }()
	if !rsm.auth.IsUserAuthed(rsm.lastIqQuery.Stanza.From) {
		return part.Channel().SendElement(RosterErrorElem(rsm.lastIqQuery.Stanza.ID, rsm.lastIqQuery.Stanza.From, "auth", ErrRMForbidden))
	}
	switch rsm.lastMatchType {
	case rosterQuery:
		return rsm.handleQuery(part)
	case rosterSet:
		return rsm.handleSet(part)
	}
	return nil
}

func (rsm *RosterManager) handleQuery(part xc.Part) error {
	rosters := []Roster{}
	rsm.handler.Rosters(&rosters)
	iqq := &rsm.lastIqQuery
	return part.Channel().SendElement(rsm.queryResult(rosters, iqq.Stanza.ID, part.Attr().Domain, xc.StanzaResult))
}

func (rsm *RosterManager) queryResult(rosters []Roster, id, to string, typ xc.StanzaType) stravaganza.Element {
	elems := []stravaganza.Element{}
	for _, roster := range rosters {
		elems = append(elems, roster.ToElem())
	}
	var elem stravaganza.Element
	IqQuery{Stanza: xc.Stanza{
		Name: xc.NameIQ,
		ID:   id,
		Type: typ,
		To:   to}, Xmlns: nsRoster, Children: elems}.ToElem(&elem)
	return elem
}

func (rsm *RosterManager) handleSet(part xc.Part) error {
	rosters := []Roster{}
	iqq := &rsm.lastIqQuery
	rsm.rostersFromQuery(&rosters)
	if len(rosters) != 1 {
		return part.Channel().SendElement(RosterErrorElem(iqq.Stanza.ID, iqq.Stanza.From, "modify", ErrRMBadRequest))
	}
	roster := rosters[0]
	var err error
	switch roster.Sub {
	case SubRm:
		err = rsm.handler.Delete(roster)
	default:
		err = rsm.handler.AddOrUpdate(roster)
	}
	if err != nil {
		part.Logger().Printf(xc.Error, "set roster: %s\n", err.Error())
		return part.Channel().SendElement(RosterErrorElem(iqq.Stanza.ID, iqq.Stanza.From, "modify", ErrRMNotAcceptable))
	}
	return part.Channel().SendElement(xc.Stanza{
		Name: xc.NameIQ,
		ID:   iqq.Stanza.ID,
		Type: xc.StanzaResult,
		To:   iqq.Stanza.From}.ToElemBuilder().Build())
}

func (rsm *RosterManager) rostersFromQuery(rosters *[]Roster) {
	for _, item := range rsm.lastIqQuery.Children {
		*rosters = append(*rosters, RosterFromElem(item))
	}
}
