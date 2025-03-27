IMAGE_NAME=tg-bot
CONTAINER_NAME=tg-bot-container
ENV_FILE=.env

build:
	@echo "Building the Docker image..."
	docker build -t $(IMAGE_NAME) .

run:
	@echo "Running the Docker container..."
	docker run -d --name $(CONTAINER_NAME) --restart unless-stopped --env-file $(ENV_FILE) $(IMAGE_NAME)

stop:
	@echo "Stopping the container..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

restart: stop run

logs:
	docker logs -f $(CONTAINER_NAME)

clean:
	docker rmi $(IMAGE_NAME)