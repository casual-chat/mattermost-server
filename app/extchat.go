// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (a *App) IsLinked(username string, platform string) bool {
	_, err := a.Srv().Store.ExtRef().GetByRealUserIdAndPlatform(username, platform)
	return err == nil
}

func (a *App) LinkAccount(extRef *model.ExtRef) *model.AppError {
	_, err := a.Srv().Store.ExtRef().Save(extRef)
	if err != nil {
		return model.NewAppError("LinkAccount", "app.ext_ref.link_account.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// func (a *App) CreateAliasAccount(userName string, externalId string, platform string) *model.AppError {
// 	userModel := &model.User{Email: "",
// 		Nickname: userName,
// 		Password: "",
// 		Username: userName,
// 		IsAlias:  true, //this is a computed field
// 	}
// 	user, err := a.Srv().Store.User().Save(userModel)
// 	if err != nil {
// 		return model.NewAppError("CreateAlias", "app.ext_ref.create_alias.internal_error", nil, err.Error(), http.StatusInternalServerError)
// 	}
// 	ext_ref := &model.ExtRef{
// 		RealUserId:       "",
// 		ExternalId:       externalId,
// 		ExternalPlatform: platform,
// 		AliasUserId:      user.Id,
// 	}
// 	_, extRefErr := a.Srv().Store.ExtRef().Save(ext_ref)
// 	if extRefErr != nil {
// 		return model.NewAppError("CreateAlias", "app.ext_ref.save_ext_ref.internal_error", nil, err.Error(), http.StatusInternalServerError)
// 	}

// 	return nil
// }

func (a *App) GetOrCreateAliasUserId(userName string, externalId string, platform string) (string, *model.AppError) {
	ext_ref, err := a.Srv().Store.ExtRef().GetByExtIdAndPlatform(externalId, platform)
	if err != nil {
		ext_ref = &model.ExtRef{
			RealUserId:       "",
			ExternalId:       externalId,
			ExternalPlatform: platform,
			AliasUserId:      "",
		}
		//return "", model.NewAppError("GetAlias", "app.ext_ref.get_alias.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	if ext_ref.AliasUserId != "" {
		return ext_ref.AliasUserId, nil
	}
	nameWithPlatform := fmt.Sprintf("%s (%s)", userName, platform)
	nameWithoutSpaces := strings.ReplaceAll(userName, " ", "-")
	userModel := &model.User{
		Email:    strings.ToLower(fmt.Sprintf("%s@%s.com", nameWithoutSpaces, platform)),
		Nickname: nameWithPlatform,
		Password: "",
		Username: strings.ToLower(nameWithoutSpaces),
		IsAlias:  true,
	}
	user, usr_err := a.Srv().Store.User().Save(userModel)
	if usr_err != nil {
		return "", model.NewAppError("CreateAlias", "app.ext_ref.create_alias.internal_error", nil, usr_err.Error(), http.StatusInternalServerError)
	}
	ext_ref.AliasUserId = user.Id
	_, save_err := a.Srv().Store.ExtRef().Save(ext_ref)
	if save_err != nil {
		return "", model.NewAppError("SaveExtRef", "app.ext_ref.save_ext_ref.internal_error", nil, save_err.Error(), http.StatusInternalServerError)
	}

	return user.Id, nil
}

func (a *App) GetExtRefFromAliasUserId(aliasId string) (*model.ExtRef, *model.AppError) {
	ext_ref, err := a.Srv().Store.ExtRef().GetByAliasUserId(aliasId)
	if err != nil {
		return nil, model.NewAppError("GetAlias", "app.ext_ref.get_alias.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return ext_ref, nil
}
