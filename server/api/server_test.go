package api

import (
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/sordgom/PasswordManager/server/config"
	"github.com/sordgom/PasswordManager/server/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type MockAppContext struct {
}

func newTestServer(t *testing.T, vaultService config.VaultService) *Server {
	config := config.Config{
		RedisAddress: "localhost:6379",
	}

	server, err := NewServer(config, vaultService)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func randomVault() *model.Vault {
	return &model.Vault{
		Name:           RandomString(10),
		MasterPassword: RandomString(10),
		Passwords:      []model.Password{},
	}
}

func RandomPassword(v *model.Vault) {
	v.NewPassword(RandomString(10), RandomString(10), RandomString(10), RandomString(10), RandomString(10))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) //0~max-min
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] //0~k-1
		sb.WriteByte(c)
	}

	return sb.String()
}
