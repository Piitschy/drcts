run:
	@echo "Running go run ."
	go run .

install:
	@echo "Installing..."
	go mod tidy
	go install github.com/goreleaser/goreleaser/v2@latest
	@echo "Done!"

release:
	@if [ "$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then echo "You are not on 'main' branch"; exit 1; fi
	@echo "Set the version first with 'git tag -a v0.0.0 -m 'release note''"
	@read -p "Press enter to continue"
	@echo "Releasing..."
	sh ./update-version.sh
	goreleaser release --clean
	@echo "Done!"
