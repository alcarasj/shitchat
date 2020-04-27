.PHONY: build directory peer peers-macos

build:
	go build -o ./bin/directory -v ./directory
	go build -o ./bin/peer -v ./peer

directory: 
	go run directory/directory.go

peer:
	go run peer/peer.go -username=peerone -port=3000

peers-macos:
	osascript -e 'tell app "Terminal" to do script "cd ~/Code/shitchat && go run peer.go -username=peerone -port=3000"'
	osascript -e 'tell app "Terminal" to do script "cd ~/Code/shitchat && go run peer.go -username=peertwo -port=3001"'
	osascript -e 'tell app "Terminal" to do script "cd ~/Code/shitchat && go run peer.go -username=peerthree -port=3002"'
