package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sordgom/PasswordManager/server/mocks"
	"github.com/sordgom/PasswordManager/server/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreatePassword(t *testing.T) {
	// Create a test vault and add a random password to it
	vault := randomVault()
	name := RandomString(10)
	url := RandomString(10)
	username := RandomString(10)
	password := RandomString(10)
	hint := RandomString(10)

	testCases := []struct {
		name          string
		body          gin.H
		vaultName     string
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"name":     name,
				"url":      url,
				"username": username,
				"password": password,
				"hint":     hint,
			},
			vaultName: vault.Name,
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(vault, nil)
				mock.EXPECT().SaveVaultToRedis(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				require.Len(t, vault.Passwords, 1)
				require.Equal(t, `{"message":"Password added successfully"}`, recorder.Body.String())
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

			url := fmt.Sprintf("/password?vault_name=%s", tc.vaultName)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetPasswords(t *testing.T) {
	// Generate a new vault with 2 random passwords
	vault := randomVault()
	RandomPassword(vault)
	RandomPassword(vault)

	testCases := []struct {
		name          string
		vaultName     string
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			vaultName: vault.Name,
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(vault, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPasswords(t, recorder.Body, vault.Passwords)
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

			url := fmt.Sprintf("/passwords?vault_name=%s", tc.vaultName)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetPassword(t *testing.T) {
	// Generate a new vault with 2 random passwords
	vault := randomVault()
	password := vault.NewPassword("name1", "url1", "username1", "password1", "hint1")

	testCases := []struct {
		name          string
		vaultName     string
		passwordName  string
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			vaultName:    vault.Name,
			passwordName: "name1",
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(vault, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPassword(t, recorder.Body, password)
				require.Equal(t, password, vault.Passwords[0])
			},
		},
		{
			name:         "Password Name Not Found",
			vaultName:    vault.Name,
			passwordName: "name2",
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(vault, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, `{"error":"Password not found"}`, recorder.Body.String())
			},
		},
		{
			name:         "Vault Not Found",
			vaultName:    "vault2",
			passwordName: "name1",
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().LoadVaultFromRedis(gomock.Any(), gomock.Any()).Times(1).Return(nil, fmt.Errorf("Vault not found"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Equal(t, `{"error":"Vault not found"}`, recorder.Body.String())
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

			url := fmt.Sprintf("/password?vault_name=%s&password_name=%s", tc.vaultName, tc.passwordName)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchPasswords(t *testing.T, body *bytes.Buffer, passwords []model.Password) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotPasswords []model.Password
	err = json.Unmarshal(data, &gotPasswords)
	require.NoError(t, err)
	require.Equal(t, passwords, gotPasswords)
}

func requireBodyMatchPassword(t *testing.T, body *bytes.Buffer, password model.Password) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotPassword model.Password
	err = json.Unmarshal(data, &gotPassword)
	require.NoError(t, err)
	require.Equal(t, password, gotPassword)
}
