.PHONY: build frontend backend clean run

frontend:
	cd web && npm install && npm run build

backend:
	go clean -cache
	go build -ldflags="-s -w" -o gopanel .

build: frontend backend

run: build
	./gopanel --config ./gopanel.json

dev-frontend:
	cd web && npm run dev

dev-backend:
	go run . --config ./gopanel.json

clean:
	rm -f gopanel gopanel.exe
	rm -rf web/dist

install:
	./install.sh
