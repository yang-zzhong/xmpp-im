// Copyright (c) 2021 Yang,Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package xmppim

import (
	"github.com/jackal-xmpp/stravaganza/v2"
	xc "github.com/yang-zzhong/xmpp-core"
)

type PresenceClient struct{}

func (pc PresenceClient) Sub(jid string, part xc.Part) {

}

func (pc PresenceClient) Unsub(jid string, part xc.Part) {

}

func (pc *PresenceClient) Match(elem stravaganza.Element) bool {
	return true
}

func (pc *PresenceClient) Handle(_ stravaganza.Element, part xc.Part) error {
	return nil
}
