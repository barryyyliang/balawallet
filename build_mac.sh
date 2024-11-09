cd tool
rm tools
go build tools.go && ./tools && cd .. && go build && mv Web3Wallet bin/
