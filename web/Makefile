WEBPACK_IMAGE = node:23-alpine
PYTHON_IMAGE = python-requests

install-npm:
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npm install
build-frontend:
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npx webpack
build-python:
	docker build -t $(PYTHON_IMAGE) .
update-repo-versions:
	docker run --rm $(if ${GITHUB_TOKEN},-e GITHUB_TOKEN) -v $(PWD):/app $(PYTHON_IMAGE) python /app/update_repo.py

update-pages: install-npm build-frontend build-python update-repo-versions


.PHONY: update-pages install-npm build-frontend build-python build-repo-versions

