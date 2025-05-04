WEBPACK_IMAGE = node:23-alpine
PYTHON_IMAGE = python-requests

update-pages: mkdir-dist build-python update-repo-versions get-issues install-npm cp-to-src-repo-json cp-to-src-issues-json build-frontend cp-to-dist-repo-json

install-npm:
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npm install
build-frontend:
	docker run --rm -v $(PWD):/app -w /app $(WEBPACK_IMAGE) npx webpack
cp-to-src-repo-json:
	cp repo.json src/repo.json
cp-to-dist-repo-json:
	cp repo.json dist/repo.json
cp-to-src-issues-json:
	cp issues.json src/issues.json
build-python:
	docker build -t $(PYTHON_IMAGE) .
update-repo-versions:
	docker run --rm $(if ${GITHUB_TOKEN},-e GITHUB_TOKEN) -v $(PWD):/app $(PYTHON_IMAGE) python /app/update_repo.py
get-issues:
	docker run --rm $(if ${GITHUB_TOKEN},-e GITHUB_TOKEN) -v $(PWD):/app $(PYTHON_IMAGE) python /app/issues.py
mkdir-dist:
	mkdir -p dist

.PHONY: update-pages install-npm build-frontend build-python build-repo-versions cp-to-src-repo-json cp-to-dist-repo-json cp-to-src-issues-json

