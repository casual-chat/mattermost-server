// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"testing"
	"github.com/stretchr/testify/require"

)
func TestFriendList(t *testing.T){
	th := Setup(t).InitBasic()
	defer th.TearDown()
	Client := th.Client



	_, resp := Client.SendFriendRequest(th.BasicUser2.Id,th.BasicUser.Id)
	CheckNoError(t, resp)

	list, resp := Client.GetReceviedList(th.BasicUser.Id)
	CheckNoError(t, resp)
	require.Equal(t, len(list), 1, "no pending request")

	_, resp = Client.AcceptFriendRequest(th.BasicUser2.Id, th.BasicUser.Id)
	CheckNoError(t, resp)

	list, resp = Client.GetReceviedList(th.BasicUser.Id)
	CheckNoError(t, resp)
	require.Equal(t, len(list), 0, "request not accepted")

	list, resp = Client.GetFriendList(th.BasicUser.Id)
	CheckNoError(t, resp)
	require.Equal(t, len(list), 1, "request not accepted")

	_, resp = Client.RejectFriendRequest(th.BasicUser2.Id, th.BasicUser.Id)
	CheckNoError(t, resp)

	list, resp = Client.GetReceviedList(th.BasicUser.Id)
	CheckNoError(t, resp)
	require.Equal(t, len(list), 0, "request not deleted")

	list, resp = Client.GetFriendList(th.BasicUser.Id)
	CheckNoError(t, resp)
	require.Equal(t, len(list), 0, "request not d")


	//require.Equal(t, newEmoji.Name, emoji.Name, "create with wrong name")

}