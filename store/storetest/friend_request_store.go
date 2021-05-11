package storetest

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFriendRequestStore(t *testing.T, ss store.Store) {

	t.Run("Save", func(t *testing.T) { testSave(t, ss) })
	// t.Run("GetByAliasUserId", func(t *testing.T) { testGetByAliasUserId(t, ss) })
	// t.Run("UpdateRealId", func(t *testing.T) { testUpdateRealId(t, ss) })
	// t.Run("GetByExtIdAndPlatform", func(t *testing.T) { testGetByExtIdAndPlatform(t, ss) })
	// t.Run("GetByRealUserIdAndPlatform", func(t *testing.T) { testGetByRealUserIdAndPlatform(t, ss) })
	// t.Run("Unlink", func(t *testing.T) { testUnlink(t, ss) })
}

// var testRealUserId1 = model.NewId()
// var testAliasUserId1 = model.NewId()
// var testExternalId1 = model.NewId()
// var testExternalPlatform1 = "Telegram"
var testSenderId = model.NewId()
var testReceiverId = model.NewId()

func testSave(t *testing.T, ss store.Store) {
	friendRequest := &model.FriendRequest{
		SenderId:		testSenderId,
		ReceiverId:		testReceiverId,
		Status:			"pending",
	}
	result, err := ss.FriendRequest().Save(friendRequest)
	require.Nil(t, err)
	assert.Equal(t, result.SenderId, testSenderId)
	assert.Equal(t, result.ReceiverId, testReceiverId)

	result, err = ss.FriendRequest().FindFriendRequest(testSenderId, testReceiverId)
	require.Nil(t, err)
	require.Equal(t, result.SenderId, testSenderId)
	require.Equal(t, result.ReceiverId, testReceiverId)
	require.Equal(t, "pending", result.Status)

	result, err = ss.FriendRequest().FindFriendRequest(testReceiverId, testSenderId)
	require.Nil(t, err)
	require.Nil(t, result)

	result1, err := ss.FriendRequest().GetReceivedList(testReceiverId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 1)

	result1, err = ss.FriendRequest().GetPendingList(testReceiverId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 0)

	result1, err = ss.FriendRequest().GetReceivedList(testSenderId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 0)

	result1, err = ss.FriendRequest().GetPendingList(testSenderId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 1)

	err = ss.FriendRequest().AcceptRequest(testSenderId, testReceiverId)
	require.Nil(t, err)

	result, err = ss.FriendRequest().FindFriendRequest(testSenderId, testReceiverId)
	require.Nil(t, err)
	require.Equal(t, "accepted", result.Status)

	result1, err = ss.FriendRequest().GetFriendList(testSenderId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 1)

	result1, err = ss.FriendRequest().GetFriendList(testReceiverId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 1)

	err = ss.FriendRequest().RemoveRequest(testSenderId, testReceiverId)
	require.Nil(t, err)
	result, err = ss.FriendRequest().FindFriendRequest(testSenderId, testReceiverId)
	require.Nil(t, err)
	require.Nil(t, result)

	result1, err = ss.FriendRequest().GetFriendList(testSenderId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 0)

	result1, err = ss.FriendRequest().GetFriendList(testReceiverId)
	require.Nil(t, err)
	require.Equal(t, len(result1), 0)


}

// func testGetByExtIdAndPlatform(t *testing.T, ss store.Store) {
// 	result, err := ss.ExtRef().GetByExtIdAndPlatform(testExternalId1, testExternalPlatform1)
// 	require.Nil(t, err)
// 	assert.Equal(t, result.RealUserId, testRealUserId1)
// }