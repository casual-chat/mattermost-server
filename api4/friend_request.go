// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (

	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (api *API) InitFriendRequest() {
	api.BaseRoutes.FriendRequest.Handle("", api.ApiSessionRequired(sendFriendRequest)).Methods("POST")
	api.BaseRoutes.FriendRequest.Handle("/reject", api.ApiSessionRequired(rejectFriendRequest)).Methods("POST")
	api.BaseRoutes.FriendRequest.Handle("/accept", api.ApiSessionRequired(acceptFriendRequest)).Methods("POST")
	api.BaseRoutes.FriendRequest.Handle("/receive", api.ApiSessionRequired(showReceivedFriendRequest)).Methods("GET")
	api.BaseRoutes.FriendRequest.Handle("/friendslist", api.ApiSessionRequired(showFriendsList)).Methods("GET")
}


func sendFriendRequest(c *Context, w http.ResponseWriter, r *http.Request){
	senderid := r.URL.Query().Get("senderid")
	receiverid := r.URL.Query().Get("receiverid")
	if senderid == "" {
		c.SetInvalidUrlParam("senderid")
		return
	}
	if receiverid == "" {
		c.SetInvalidUrlParam("receiverid")
		return
	}
	
	err := c.App.SendFriendRequest(senderid, receiverid)
	if err != nil {
		c.Err = err
		return
	}
	ReturnStatusOK(w)
}

func rejectFriendRequest(c *Context, w http.ResponseWriter, r *http.Request){
	senderid := r.URL.Query().Get("senderid")
	receiverid := r.URL.Query().Get("receiverid")
	if senderid == "" {
		c.SetInvalidUrlParam("senderid")
		return
	}
	if receiverid == "" {
		c.SetInvalidUrlParam("receiverid")
		return
	}
	
	err := c.App.RejectFriendRequest(senderid, receiverid)
	if err != nil {
		c.Err = err
		return
	}
	ReturnStatusOK(w)
}

func acceptFriendRequest(c *Context, w http.ResponseWriter, r *http.Request){
	senderid := r.URL.Query().Get("senderid")
	receiverid := r.URL.Query().Get("receiverid")
	if senderid == "" {
		c.SetInvalidUrlParam("senderid")
		return
	}
	if receiverid == "" {
		c.SetInvalidUrlParam("receiverid")
		return
	}
	
	err := c.App.AcceptFriendRequest(senderid, receiverid)
	if err != nil {
		c.Err = err
		return
	}
	ReturnStatusOK(w)
}

func showReceivedFriendRequest(c *Context, w http.ResponseWriter, r *http.Request){
	receiverid := c.App.Session().UserId
	if receiverid == "" {
		c.SetInvalidUrlParam("receiverid")
		return
	}
	


	listRequest, err := c.App.ShowPendingFriendRequest(receiverid)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.FriendRequestToJson(listRequest)))
}

func showFriendsList(c *Context, w http.ResponseWriter, r *http.Request){
	receiverid := c.App.Session().UserId
	if receiverid == "" {
		c.SetInvalidUrlParam("receiverid")
		return
	}
	


	listRequest, err := c.App.ShowFriendsList(receiverid)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(model.FriendRequestToJson(listRequest)))
}
