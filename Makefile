docker-build:
	docker build -t rusprofile .

docker-run:
	docker run -p 8080:8080 rusprofile
