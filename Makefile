.PHONY: build clean

build:
	docker run \
		-v $(shell pwd):/build \
		-w /build \
		-u $(shell id -u ${USER}):$(shell id -g ${USER}) \
		-e DOTNET_CLI_TELEMETRY_OPTOUT=1 \
		-e DOTNET_CLI_HOME=/tmp/.dotnet \
		mcr.microsoft.com/dotnet/sdk:5.0 dotnet publish -c release

install:
	rm /usr/local/bin/dbcore
	ln -s $(CURDIR)/bin/release/netcoreapp3.0/linux-x64/publish/dbcore /usr/local/bin
