mockVault :
	mockgen -destination=mocks/mock_vaultservice.go -package=mocks github.com/sordgom/PasswordManager/server/config VaultService
dockerDown:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml rm -f && docker rmi api:1.0
dockerUp:
	docker compose -f docker-compose.yml up -d
restart:
	make dockerDown && make dockerUp
docker-dev:
	docker-compose -f docker-compose.dev.yml up
docker-test: 
	docker-compose f docker-compose.test.yml up
.PHONY: mockVault dockerDown dockerUp restart docker-dev docker-test