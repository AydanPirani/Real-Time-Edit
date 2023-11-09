CXX = g++
CXXFLAGS = -std=c++11 -Wall
SRC_DIR = crdt
BUILD_DIR = build

# Go targets
GO_SOURCES = **.go
GO_NODE = node
GO_CLIENT = client

# The treedoc target executable
TARGET = treedoc

# List your source files
SRCS = $(wildcard $(SRC_DIR)/*.cpp)

# Generate a list of object files based on the source files
OBJS = $(patsubst $(SRC_DIR)/%.cpp, $(BUILD_DIR)/%.o, $(SRCS))


# Build the C++ executable
$(TARGET): $(OBJS)
	$(CXX) $(CXXFLAGS) -o $@ $^

# Build object files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(dir $@)
	$(CXX) $(CXXFLAGS) -c -o $@ $<

# Specify phony targets
.PHONY: all clean run

all: $(TARGET) go

go: $(GO_SOURCES)
	go build -race -o $(GO_NODE) .
	go build -o $(GO_CLIENT) ./rtclbedit_client

run: all
	bash run.sh 5

stop:
	- pkill -9 -f "./$(GO_NODE)"

clean: stop
	- rm $(GO_NODE) $(GO_CLIENT)
	- rm -rf $(BUILD_DIR) $(TARGET)
