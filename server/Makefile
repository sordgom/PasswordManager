mockVault :
	mockgen -destination=mocks/mock_vaultservice.go -package=mocks github.com/sordgom/PasswordManager/server/config VaultService
dockerDown:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml rm -f
dockerUp:
	docker compose -f docker-compose.yml up -d
.PHONY: mockVault dockerDown dockerUp