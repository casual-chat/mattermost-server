// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (a *App) SendFriendRequest(senderid string, receiverid string) *model.AppError {
	request, err := a.Srv().Store.FriendRequest().FindFriendRequest(senderid, receiverid)

	if err != nil {
		return model.NewAppError("SendFriendRequest", "app.friend_request.find_err", nil, err.Error(), http.StatusInternalServerError)
	}
	if request != nil {
		return model.NewAppError("SendFriendRequest", "app.friend_request.already_in_relation", nil, "already in relation", http.StatusInternalServerError)
	}

	friend_request := &model.FriendRequest{
		SenderId:   senderid,
		ReceiverId: receiverid,
		Status:     "pending",
	}

	_, err = a.Srv().Store.FriendRequest().Save(friend_request)
	if err != nil {
		return model.NewAppError("sendFriendRequest", "api.friend_request.sendFriendRequest", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (a *App) RejectFriendRequest(senderid string, receiverid string) *model.AppError {
	_, err := a.Srv().Store.FriendRequest().FindFriendRequest(senderid, receiverid)

	if err != nil {
		return model.NewAppError("rejectFriendRequest", "api.friend_request.rejectFriendRequest", nil, "", http.StatusBadRequest)
	}

	err = a.Srv().Store.FriendRequest().RemoveRequest(senderid, receiverid)
	if err != nil {
		return model.NewAppError("rejectFriendRequest", "api.friend_request.rejectFriendRequest", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (a *App) AcceptFriendRequest(senderid string, receiverid string) *model.AppError {
	_, err := a.Srv().Store.FriendRequest().FindFriendRequest(senderid, receiverid)

	if err != nil {
		return model.NewAppError("acceptFriendRequest", "api.friend_request.acceptFriendRequest", nil, "", http.StatusBadRequest)
	}

	err = a.Srv().Store.FriendRequest().AcceptRequest(senderid, receiverid)
	if err != nil {
		return model.NewAppError("rejectFriendRequest", "api.friend_request.acceptFriendRequest", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (a *App) ShowPendingFriendRequest(receiverid string) ([]*model.FriendRequest, *model.AppError) {
	list, err := a.Srv().Store.FriendRequest().GetReceivedList(receiverid)
	if err != nil {
		return nil, model.NewAppError("showPendingFriendRequest", "api.friend_request.showPendingFriendRequest", nil, "", http.StatusBadRequest)
	}
	return list, nil

}

func (a *App) ShowFriendsList(receiverid string) ([]*model.FriendRequest, *model.AppError) {
	list, err := a.Srv().Store.FriendRequest().GetFriendList(receiverid)
	if err != nil {
		return nil, model.NewAppError("showFriendList", "api.friend_request.showFriendList", nil, "", http.StatusBadRequest)
	}
	return list, nil

}
