TARGET = db 

SOURCES = $(wildcard src/*.c)

all: $(TARGET)

$(TARGET): $(SOURCES)
	gcc -Wall -Wextra -O2 -o $@ $(SOURCES)

clean:
	rm -f $(TARGET)
	rm -f db.txt

