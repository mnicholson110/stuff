build: ./src/Application.cc
	mkdir -p ./bin
	g++ -Wall ./src/*.cc ./src/vendor/*/*.cc -lGL -lGLEW -lglfw -o ./bin/game

run:
	./bin/game

clean:
	rm -rf ./bin

test: clean build run