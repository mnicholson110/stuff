
# Variables
CC = gcc
CFLAGS = -Wall -Wextra -O2
SRCDIR = src
TARGET = db 

# Find all .c files in the src/ directory
SOURCES = $(wildcard $(SRCDIR)/*.c)

# Default target
all: $(TARGET)

# Compile all .c files into a single binary
$(TARGET): $(SOURCES)
	$(CC) $(CFLAGS) -o $@ $(SOURCES)

# Clean up the binary
clean:
	rm -f $(TARGET)
	rm -f db.txt

# Phony targets to prevent conflicts with file names
.PHONY: all clean
