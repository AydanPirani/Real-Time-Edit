# all: **.go
# 	go build -race -o node .
# 	go build -o client ./rtclbedit_client

# run: all
# 	bash run.sh 5

# stop:
# 	pkill -9 -f "./node"

# clean: stop
# 	rm ./node ./client
	
CC = g++
CFLAGS = -std=c++11

SRCS = crdt/main.cpp

all: crdt node client

crdt: $(SRCS)
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
	rm -f crdt node client
