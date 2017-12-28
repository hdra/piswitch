default: setup
setup:
	@make build
	@make install
build:
	@GOOS=linux GOARCH=arm go build
install:
	@scp piSwitch pi@raspberrypi.local:
