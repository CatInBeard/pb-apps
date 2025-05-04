WEBPACK_IMAGE = node:23-alpine

build:
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npm install
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npx webpack

.PHONY: build
