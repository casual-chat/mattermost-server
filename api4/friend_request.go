// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/mattermost/mattermost-server/v5/model"
// )

// func (api *API) InitFriendRequest() {
// 	api.BaseRoutes.Emojis.Handle("", api.ApiSessionRequired(sendFriendRequest)).Methods("POST")
// 	api.BaseRoutes.Emojis.Handle("reject", api.ApiSessionRequired(rejectFriendRequest)).Methods("POST")
// 	api.BaseRoutes.Emojis.Handle("accept", api.ApiSessionRequired(acceptFriendRequest)).Methods("POST")
// }


// func sendFriendRequest(c *Context, w http.ResponseWriter, r *http.Request){
// 	senderid := r.URL.Query().Get("senderid")
// 	receiverid := r.URL.Query().Get("receiverid")
// 	if senderid == "" {
// 		c.SetInvalidUrlParam("senderid")
// 		return
// 	}
// 	if receiverid == "" {
// 		c.SetInvalidUrlParam("receiverid")
// 		return
// 	}
	
// 	err := c.App.SendFriendRequest(senderid, receiverid)
// 	if err != nil {
// 		c.Err = err
// 		return
// 	}
// 	ReturnStatusOK(w)
// }
