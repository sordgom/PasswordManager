package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sordgom/PasswordManager/server/mocks"
	"github.com/sordgom/PasswordManager/server/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateVault(t *testing.T) {
	name := RandomString(10)
	masterPassword := RandomString(10)
	vault := model.Vault{
		Name:           name,
		MasterPassword: masterPassword,
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"name":            name,
				"master_password": masterPassword,
			},
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
				mock.EXPECT().SaveVaultToRedis(gomock.Any(), gomock.Eq(&vault)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				require.Equal(t, `{"message":"Vault created successfully"}`, recorder.Body.String())
				require.Equal(t, vault.Name, name)
				require.Equal(t, vault.MasterPassword, masterPassword)
			},
		},
		{
			name: "Empty Name",
			body: gin.H{
				"name":            "",
				"master_password": masterPassword,
			},
			buildStubs: func(mock *mocks.MockVaultService) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, `{"message":"Name and master password are required"}`, recorder.Body.String())
			},
		},
		{
			name: "Empty MP",
			body: gin.H{
				"name":            name,
				"master_password": "",
			},
			buildStubs: func(mock *mocks.MockVaultService) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, `{"message":"Name and master password are required"}`, recorder.Body.String())
			},
		},
		// {
		// 	name: "Empty MP",
		// 	body: gin.H{
		// 		"name":            name,
		// 		"master_password": "",
		// 	},
		// 	vaultName:      name,
		// 	masterPassword: masterPassword,
		// 	buildStubs: func(mock *mocks.MockVaultService) {
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 		require.Equal(t, `{"message":"Name and master password are required"}`, recorder.Body.String())
		// 	},
		// },
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockVaultService(ctrl)
			tc.buildStubs(mock)

			server := newTestServer(t, mock)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/vault"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestVerifyMasterPassword(t *testing.T) {
	vaultName := RandomString(10)
	masterPassword := RandomString(10)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Verified",
			body: gin.H{
				"vault_name":      vaultName,
				"master_password": masterPassword,
			},
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().VerifyMasterPassword(gomock.Any(), vaultName, masterPassword).Times(1).Return(true)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, `{"message":"Master password is verified"}`, recorder.Body.String())
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"vault_name":      vaultName,
				"master_password": "randompassword",
			},
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().VerifyMasterPassword(gomock.Any(), vaultName, "randompassword").Times(1).Return(false)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				require.Equal(t, `{"message":"Master password is not verified"}`, recorder.Body.String())
			},
		},
		{
			name: "Empty Vault Name",
			body: gin.H{
				"vault_name":      "",
				"master_password": masterPassword,
			},
			buildStubs: func(mock *mocks.MockVaultService) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, `{"message":"Vault name is required"}`, recorder.Body.String())
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockVaultService(ctrl)
			tc.buildStubs(mock)

			server := newTestServer(t, mock)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/vault/verify"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
