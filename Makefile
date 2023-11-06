all: **.go
	go build -o node .
	go build -o client ./rtclbedit_client
