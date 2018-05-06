default: setup
setup:
	@make build.client
	@make build.server
	@make install
build.client:
	@cd client && yarn install
	@cd client && yarn build
	@statik -src=client/build
build.server:
	@GOOS=linux GOARCH=arm go build -o build/piswitch
install:
	@scp build/piswitch pi@raspberrypi.local:piswitch/piswitch
