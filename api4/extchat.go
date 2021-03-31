// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (api *API) InitExtChat() {
	api.BaseRoutes.ExtChat.Handle("/isLinked", api.ApiHandler(isLinked)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/linkAccount", api.ApiSessionRequired(linkAccount)).Methods("POST")
	api.BaseRoutes.ExtChat.Handle("/post", api.ApiHandler(postToChannel)).Methods("POST")
	api.BaseRoutes.ExtChat.Handle("/aliasUserId", api.ApiHandler(getAliasUserId)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/ref", api.ApiHandler(getExtRef)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/refByChannel", api.ApiHandler(getExtRefByChannelId)).Methods("GET")
	api.BaseRoutes.ExtChat.Handle("/channel", api.ApiHandler(getExtChannelId)).Methods("GET")

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

func getAliasUserId(c *Context, w http.ResponseWriter, r *http.Request) {
	externalPlatform := c.Params.ExtChatPlatform
	externalId := r.URL.Query().Get("externalId")
	username := r.URL.Query().Get("username")
	userId, err := c.App.GetOrCreateAliasUserId(username, externalId, externalPlatform)
	if err != nil {
		c.Err = err
		return
	}
	json_str, json_err := json.Marshal(userId)
	if json_err != nil {
		c.Err = model.NewAppError("AliasToJson", "app.ext_ref.create_alias.internal_error", nil, json_err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json_str)
}

func getExtRef(c *Context, w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	ext_ref, err := c.App.GetExtRefFromAliasUserId(userId)
	if err != nil {
		w.Write([]byte("{}"))
		return
	}

	w.Write([]byte(ext_ref.ToJson()))
}

func getExtRefByChannelId(c *Context, w http.ResponseWriter, r *http.Request) {
	channelId := r.URL.Query().Get("channelId")
	channel, err := c.App.GetChannel(channelId)
	if err != nil {
		c.Err = err
		return
	}
	if channel.Type != model.CHANNEL_DIRECT {
		c.Err = model.NewAppError("GetExtRefByChannel", "app.ext_ref.get_by_channel.internal_error", nil, "Channel not direct", http.StatusInternalServerError)
		return
	}
	members, members_err := c.App.GetChannelMembersPage(channelId, 0, 2)
	if members_err != nil {
		c.Err = members_err
		return
	}
	var aliasId string = ""

	for _, member := range *members {
		_, ext_err := c.App.GetExtRefFromAliasUserId(member.UserId)

		if ext_err == nil {
			aliasId = member.UserId
			break
		}

	}

	if aliasId == "" {
		w.Write([]byte("{}"))
		return
	}
	ext_ref, ext_ref_err := c.App.GetExtRefFromAliasUserId(aliasId)
	if ext_ref_err != nil {
		w.Write([]byte("{}"))
		return
	}

	w.Write([]byte(ext_ref.ToJson()))
}

func postToChannel(c *Context, w http.ResponseWriter, r *http.Request) {
	post := model.PostFromJson(r.Body)
	userId := post.UserId
	_, err := c.App.GetUser(userId)
	if err != nil {
		c.Err = err
		return
	}

	_, err_create := c.App.CreatePostAsUser(post, userId, false)
	if err_create != nil {
		c.Err = err_create
		return
	}
	ReturnStatusOK(w)
}

func getExtChannelId(c *Context, w http.ResponseWriter, r *http.Request) {
	externalPlatform := c.Params.ExtChatPlatform
	externalId := r.URL.Query().Get("externalId")
	userId := c.App.Session().UserId
	aliasId, err := c.App.GetAliasUserId(externalId, externalPlatform)
	if err != nil {
		c.Err = err
		return
	}
	channelId, channel_err := c.App.GetExtChannelIdByUsers(userId, aliasId)
	if channel_err != nil {
		c.Err = channel_err
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"channelId\":\"%s\"}", channelId)))
}
