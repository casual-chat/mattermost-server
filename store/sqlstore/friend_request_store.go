// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"
	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/pkg/errors"
)

type SqlFriendRequestStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

func newSqlFriendRequestStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.FriendRequestStore {
	s := &SqlFriendRequestStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.FriendRequest{}, "FriendRequest").SetKeys(false, "SenderId", "ReceiverId")
		table.ColMap("SenderId").SetMaxSize(64)
		table.ColMap("ReceiverId").SetMaxSize(64)
		table.ColMap("Status").SetMaxSize(64)
	}

	return s
}

func (es SqlFriendRequestStore) createIndexesIfNotExists() {
	es.CreateIndexIfNotExists("idx_friend_request", "FriendRequest", "SenderId")
	es.CreateIndexIfNotExists("idx_friend_request", "FriendRequest", "ReceiverId")
	es.CreateIndexIfNotExists("idx_friend_request", "FriendRequest", "Status")
}



func (es SqlFriendRequestStore) Save(request *model.FriendRequest) (*model.FriendRequest, error) {

	if err := es.GetMaster().Insert(request); err != nil {
		return nil, errors.Wrap(err, "error saving request")
	}

	return request, nil
}

func (es SqlFriendRequestStore) FindFriendRequest(senderId string, receiverId string) (*model.FriendRequest, error)  {
	var friend_request *model.FriendRequest

	err := es.GetReplica().SelectOne(&friend_request,
		`SELECT
			*
		FROM
			FriendRequest
		WHERE
			SenderId = :Key1
			AND ReceiverId = :Key2`, map[string]string{"Key1": senderId, "Key2": receiverId})
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Could not find friend request")
	}

	return friend_request,nil
}

func (es SqlFriendRequestStore) GetPendingList(senderId string) ([]*model.FriendRequest, error) {
	var friend_requests []*model.FriendRequest
	var status = "pending"

	_, err := es.GetReplica().Select(&friend_requests,
		`SELECT
			*
		FROM
			FriendRequest
		WHERE
			SenderId = :Key1
			AND Status = :Key2`, map[string]string{"Key1": senderId, "Key2": status})
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of pending friend requests")
	}
	return friend_requests,nil
}

func (es SqlFriendRequestStore) GetReceivedList(receiverId string) ([]*model.FriendRequest, error) {
	var friend_requests []*model.FriendRequest
	var status = "pending"

	_, err := es.GetReplica().Select(&friend_requests,
		`SELECT
			*
		FROM
			FriendRequest
		WHERE
			ReceiverId = :Key1
			AND Status = :Key2`, map[string]string{"Key1": receiverId, "Key2": status})
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of pending friend requests")
	}
	return friend_requests, nil
}

func (es SqlFriendRequestStore) GetFriendList(userid string) ([]*model.FriendRequest, error) {
	var friend_requests []*model.FriendRequest
	var status = "accepted"

	_, err := es.GetReplica().Select(&friend_requests,
		`SELECT
			*
		FROM
			FriendRequest
		WHERE
			(ReceiverId = :Key1 OR  SenderId = :Key1)
			AND Status = :Key2`, map[string]string{"Key1": userid, "Key2": status})
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of friends")
	}
	return friend_requests, nil
}
// func (es SqlEmojiStore) GetList(offset, limit int, sort string) ([]*model.Emoji, error) {
// 	var emoji []*model.Emoji

// 	query := "SELECT * FROM Emoji WHERE DeleteAt = 0"

// 	if sort == model.EMOJI_SORT_BY_NAME {
// 		query += " ORDER BY Name"
// 	}

// 	query += " LIMIT :Limit OFFSET :Offset"

// 	if _, err := es.GetReplica().Select(&emoji, query, map[string]interface{}{"Offset": offset, "Limit": limit}); err != nil {
// 		return nil, errors.Wrap(err, "could not get list of emojis")
// 	}
// 	return emoji, nil
// }

func (es SqlFriendRequestStore) RemoveRequest(senderId string, receiverId string) error {
	sql := `DELETE
		FROM
			FriendRequest
		WHERE
			SenderId = :senderId
			AND ReceiverId = :receiverId
			AND Status = :status`

	queryParams := map[string]string{
		"senderId":       senderId,
		"receiverId": receiverId,
		"status": "accepted",
	}
	_, err := es.GetMaster().Exec(sql, queryParams)
	if err != nil {
		return err
	}
	return nil
}

func (es SqlFriendRequestStore) AcceptRequest(senderId string, receiverId string) error {
	if sqlResult, err := es.GetMaster().Exec(
		`UPDATE
			FriendRequest
		SET
			Status = :status
		WHERE
			SenderId = :senderId
			AND ReceiverId = :receiverId`, 
			map[string]interface{}{"status": "accepted", 
			"senderId": senderId, "receiverId": receiverId}); err != nil {
		return errors.Wrap(err, "could not accept request")
	} else if rows, _ := sqlResult.RowsAffected(); rows == 0 {
		return store.NewErrNotFound("FriendRequest", receiverId)
	}
 	return nil
}