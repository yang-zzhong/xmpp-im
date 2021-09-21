// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import "github.com/jackal-xmpp/stravaganza/v2"

const (
	SubNone = sub("none")
	SubTo   = sub("to")
	SubFrom = sub("from")
	SubBoth = sub("both")
	SubRm   = sub("remove")
)

type Roster struct {
	JID   string
	Sub   sub
	Name  string
	Group string
}

func (r *Roster) FromElem(elem stravaganza.Element) {
	r.JID = elem.Attribute("jid")
	s := sub(elem.Attribute("subscription"))
	if s != SubTo && s != SubFrom && s != SubBoth {
		s = SubNone
	}
	r.Sub = s
	r.Name = elem.Attribute("name")
	if elem.Child("group") != nil {
		r.Group = elem.Child("group").Text()
	}
}

func (r *Roster) ToElem() stravaganza.Element {
	b := stravaganza.NewBuilder("item").
		WithAttribute("jid", r.JID)
	if r.Sub != SubNone {
		b.WithAttribute("subscription", string(r.Sub))
	}
	if r.Name != "" {
		b.WithAttribute("name", r.Name)
	}
	if r.Group != "" {
		b.WithChild(stravaganza.NewBuilder("group").WithText(r.Group).Build())
	}
	return b.Build()
}

func RosterFromElem(elem stravaganza.Element) Roster {
	var roster Roster
	roster.FromElem(elem)
	return roster
}

func ElemFromRoster(roster Roster) stravaganza.Element {
	return roster.ToElem()
}
