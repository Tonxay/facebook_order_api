// docker buiold .
//docker tag go-api-web tonxay/go-api-web create tag for deploy images
docker build -t go-fiber-docker-api .
docker tag go-fiber-docker-api:latest your-dockerhub-username/go-fiber-docker-api:latest
docker login
docker push your-dockerhub-username/go-fiber-docker-api:latest
GOOS=linux GOARCH=amd64 go build -o server ./cmd/heroku
docker build -t tonxay/gowebhook:latest .
docker build --no-cache -t tonxay/gowebhook:latest .
