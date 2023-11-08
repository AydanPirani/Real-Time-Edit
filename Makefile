all: **.go
	go build -o node .
	go build -o client ./rtclbedit_client

run: all
	bash run.sh 5

stop:
	pkill -9 -f "./node"

clean: stop
	rm ./node ./client
	
