// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package mmpermissions

import (
	"testing"

	"github.com/mattermost/focalboard/server/model"
	mmpermissionsMocks "github.com/mattermost/focalboard/server/services/permissions/mmpermissions/mocks"
	permissionsMocks "github.com/mattermost/focalboard/server/services/permissions/mocks"

	mmModel "github.com/mattermost/mattermost-server/v6/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type TestHelper struct {
	t           *testing.T
	ctrl        *gomock.Controller
	store       *permissionsMocks.MockStore
	api         *mmpermissionsMocks.MockAPI
	permissions *Service
}

func SetupTestHelper(t *testing.T) *TestHelper {
	ctrl := gomock.NewController(t)
	mockStore := permissionsMocks.NewMockStore(ctrl)
	mockAPI := mmpermissionsMocks.NewMockAPI(ctrl)

	return &TestHelper{
		t:           t,
		ctrl:        ctrl,
		store:       mockStore,
		api:         mockAPI,
		permissions: New(mockStore, mockAPI),
	}
}

func (th *TestHelper) checkBoardPermissions(roleName string, member *model.BoardMember, teamID string, hasPermissionTo, hasNotPermissionTo []*mmModel.Permission) {
	for _, p := range hasPermissionTo {
		th.t.Run(roleName+" "+p.Id, func(t *testing.T) {
			th.store.EXPECT().
				GetBoard(member.BoardID).
				Return(&model.Board{ID: member.BoardID, TeamID: teamID}, nil).
				Times(1)

			th.api.EXPECT().
				HasPermissionToTeam(member.UserID, teamID, model.PermissionViewTeam).
				Return(true).
				Times(1)

			th.store.EXPECT().
				GetMemberForBoard(member.BoardID, member.UserID).
				Return(member, nil).
				Times(1)

			hasPermission := th.permissions.HasPermissionToBoard(member.UserID, member.BoardID, p)
			assert.True(t, hasPermission)
		})
	}

	for _, p := range hasNotPermissionTo {
		th.t.Run(roleName+" "+p.Id, func(t *testing.T) {
			th.store.EXPECT().
				GetBoard(member.BoardID).
				Return(&model.Board{ID: member.BoardID, TeamID: teamID}, nil).
				Times(1)

			th.api.EXPECT().
				HasPermissionToTeam(member.UserID, teamID, model.PermissionViewTeam).
				Return(true).
				Times(1)

			th.store.EXPECT().
				GetMemberForBoard(member.BoardID, member.UserID).
				Return(member, nil).
				Times(1)

			hasPermission := th.permissions.HasPermissionToBoard(member.UserID, member.BoardID, p)
			assert.False(t, hasPermission)
		})
	}
}
