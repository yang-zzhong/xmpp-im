// Copyright (c) 2021 Yang Zhong
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"net"
	"os"
	"time"

	xc "github.com/yang-zzhong/xmpp-core"
	xi "github.com/yang-zzhong/xmpp-im"
)

var (
	client        *xc.ClientPart
	server        *xc.XPart
	clientRoster  *xi.RosterClient
	rosterManager *xi.RosterManager
)

func main() {
	logger := xc.NewLogger(os.Stdout)
	go func() {
		time.Sleep(time.Second)
		conn, err := net.Dial("tcp", "127.0.0.1:4444")
		if err != nil {
			panic(err)
		}
		client = xc.NewClientPart(xc.NewTcpConn(conn, true), logger, &xc.PartAttr{})
		client.Channel().SetLogger(logger)
		clientRoster = xi.NewRosterClient("test@127.0.0.1:4444/test", &rosterHandler{})
		client.WithElemHandler(clientRoster)
		if err := client.Negotiate(); err != nil {
			panic(err)
		}
		errChan := client.Run()
		clientRoster.Query(client.Channel())
		panic(<-errChan)
	}()

	l, err := net.Listen("tcp", ":4444")
	if err != nil {
		fmt.Printf("listen: %s\n", err.Error())
		return
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("accept: %s\n", err.Error())
		return
	}
	server = xc.NewXPart(xc.NewTcpConn(conn, false), "127.0.0.1", logger)
	server.Channel().SetLogger(logger)
	rosterManager = xi.NewRosterManager(newRosterManageHandler(), &isAuthed{})
	server.WithElemHandler(rosterManager)
	panic(<-server.Run())
}

type isAuthed struct{}

func (ia *isAuthed) IsUserAuthed(jid string) bool {
	return true
}

type rosterManageHandler struct {
	rosters []xi.Roster
}

func newRosterManageHandler() *rosterManageHandler {
	return &rosterManageHandler{[]xi.Roster{
		{Name: "hello", JID: "hello@world.com", Sub: xi.SubNone, Group: "hello"},
		{Name: "hello1", JID: "hello1@world.com", Sub: xi.SubFrom, Group: "hello"},
		{Name: "hello2", JID: "hello2@world.com", Sub: xi.SubTo, Group: "hello"},
		{Name: "hello3", JID: "hello2@world.com", Sub: xi.SubBoth, Group: "hello"},
	}}
}

func (rmh *rosterManageHandler) Rosters(rosters *[]xi.Roster) error {
	*rosters = rmh.rosters
	return nil
}

func (rmh *rosterManageHandler) AddOrUpdate(roster xi.Roster) error {
	for i := range rmh.rosters {
		if rmh.rosters[i].JID == roster.JID {
			rmh.rosters[i] = roster
			return nil
		}
	}
	rmh.rosters = append(rmh.rosters, roster)
	return nil
}

func (rmh *rosterManageHandler) Delete(roster xi.Roster) error {
	for i := range rmh.rosters {
		if rmh.rosters[i].JID == roster.JID {
			rmh.rosters = append(rmh.rosters[:i], rmh.rosters[i+1:]...)
			return nil
		}
	}
	return nil
}

type rosterHandler struct{}

func (rh *rosterHandler) HandleRosterQueryResult(rosters []xi.Roster) {
	fmt.Printf("rosters in query result: %v\n", rosters)
}

func (rh *rosterHandler) HandleRosterPush(rosters []xi.Roster) {
	fmt.Printf("rosters in push result: %v\n", rosters)
}
