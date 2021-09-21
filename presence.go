// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import (
	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

type PresenceManager struct {
	lastMatchedPresence xc.Stanza
}

func (pm *PresenceManager) Match(elem stravaganza.Element) bool {
	if err := pm.lastMatchedPresence.FromElem(elem, xc.NamePresence); err != nil {
		return false
	}
	return pm.lastMatchedPresence.Name != xc.NamePresence
}

func (pm *PresenceManager) Handle(_ stravaganza.Element, part xc.Part) error {
	return nil
}
