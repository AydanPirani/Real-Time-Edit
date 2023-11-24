CC = g++
CFLAGS = -std=c++11
SRCS = crdt/main.cpp

all: treedoc node client

treedoc: $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o treedoc

node:
	go build -race -o node .

client:
	go build -o client ./rtclbedit_client

run: all
	bash run.sh 5

stop:
	pkill -9 -f "./node"

clean: stop
	rm -f treedoc node client
