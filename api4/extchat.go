// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (api *API) InitExtChat() {
	api.BaseRoutes.ExtChat.Handle("/isLinked", api.ApiHandler(isLinked)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/linkAccount", api.ApiSessionRequired(linkAccount)).Methods("POST")
	//api.BaseRoutes.ExtChat.Handle("/createAliasAccount", api.ApiHandler(createAliasAccount)).Methods("POST")
	api.BaseRoutes.ExtChat.Handle("/aliasUserId", api.ApiHandler(getAliasUserId)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/ref", api.ApiHandler(getExtRef)).Methods("GET")
}

func isLinked(c *Context, w http.ResponseWriter, r *http.Request) {
	externalPlatform := c.Params.ExtChatPlatform
	realUserId := r.URL.Query().Get("realUserId")
	isLinked := c.App.IsLinked(realUserId, externalPlatform)

	w.Write([]byte(fmt.Sprintf("%t", isLinked)))
}

func linkAccount(c *Context, w http.ResponseWriter, r *http.Request) {
	externalPlatform := c.Params.ExtChatPlatform
	realUserId := c.App.Session().UserId
	externalId := r.URL.Query().Get("externalId")
	ext_ref := &model.ExtRef{
		RealUserId:       realUserId,
		ExternalId:       externalId,
		ExternalPlatform: externalPlatform,
		AliasUserId:      "",
	}
	err := c.App.LinkAccount(ext_ref)
	if err != nil {
		c.Err = err
		return
	}
	ReturnStatusOK(w)
}

// func createAliasAccount(c *Context, w http.ResponseWriter, r *http.Request) {
// 	externalPlatform := c.Params.ExtChatPlatform
// 	externalId := r.URL.Query().Get("externalId")
// 	username := r.URL.Query().Get("nickName")
// 	err := c.App.CreateAliasAccount(username, externalId, externalPlatform)
// 	if err != nil {
// 		c.Err = err
// 		return
// 	}
// 	ReturnStatusOK(w)
// }

func getAliasUserId(c *Context, w http.ResponseWriter, r *http.Request) {
	externalPlatform := c.Params.ExtChatPlatform
	externalId := r.URL.Query().Get("externalId")
	username := r.URL.Query().Get("username")
	userId, err := c.App.GetOrCreateAliasUserId(username, externalId, externalPlatform)
	if err != nil {
		c.Err = err
		return
	}
	w.Write([]byte(userId))
}

func getExtRef(c *Context, w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	ext_ref, err := c.App.GetExtRefFromAliasUserId(userId)
	if err != nil {
		w.Write([]byte(nil))
		return
	}

	w.Write([]byte(ext_ref.ToJson()))
}
