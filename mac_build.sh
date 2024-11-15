cd tool
rm tools
go build tools.go && ./tools && cd .. && go build && mv w3wallet bin/
