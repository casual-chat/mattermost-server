// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
)

type FriendRequest struct {
	SenderId string `json:"sender_id"`
	ReceiverId  string `json:"receiver_id"`
	Status  string `json:"status"`
}

func (friend_request *FriendRequest) IsValid() *AppError {
	if !IsValidId(friend_request.SenderId) || !IsValidId(friend_request.ReceiverId) {
		return NewAppError("Friend.IsValid", "model.friend.id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(friend_request.SenderId) > 26 || len(friend_request.ReceiverId) > 26 {
		return NewAppError("Emoji.IsValid", "model.friend.user_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (friend_request *FriendRequest) ToJson() string {
	b, _ := json.Marshal(friend_request)
	return string(b)
}

func FriendRequestFromJson(data io.Reader) *FriendRequest {
	var friend_request *FriendRequest
	json.NewDecoder(data).Decode(&friend_request)
	return friend_request
}