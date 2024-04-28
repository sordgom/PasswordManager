package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sordgom/PasswordManager/server/mocks"

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
	vault.NewPassword(name, url, username, password, hint)
	vaultPassword := vault.Passwords[0]

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(mock *mocks.MockVaultService)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":     name,
				"url":      url,
				"username": username,
				"password": password,
				"hint":     hint,
			},
			buildStubs: func(mock *mocks.MockVaultService) {
				mock.EXPECT().SaveVaultToRedis(gomock.Any()).Do(func() {
					require.Len(t, vault.Passwords, 1)
					require.Equal(t, vaultPassword, vault.Passwords[0])
				})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := newTestServer(t)
			server.VaultService.Vault = vault
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/password"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
