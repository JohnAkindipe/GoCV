.PHONY: help ssh-start upload-server enter-server run-server run-client clean

# Show available commands (default target)
help:
	@echo "Available commands:"
	@echo "  make ssh-start      - Start SSH agent and add private key"
	@echo "  make upload-server  - Upload server code to remote server"
	@echo "  make enter-server   - SSH into server and navigate to practice folder"
	@echo "  make run-server     - Run server on remote via SSH"
	@echo "  make run-client     - Run client locally"
	@echo "  make clean          - Clean up generated image files"

# Start SSH agent and add key
ssh-start:
	@eval $$(ssh-agent -s) && ssh-add ~/.ssh/id_rsa_greenlight

# Upload server code to remote
upload-server:
	scp ./tcpserver.go root@161.35.36.3:/root/practice/main.go
	scp ./handlers.go root@161.35.36.3:/root/practice/handlers.go
	scp ./middleware.go root@161.35.36.3:/root/practice/middleware.go
	scp ./routes.go root@161.35.36.3:/root/practice/routes.go
	scp ./webrtc.go root@161.35.36.3:/root/practice/webrtc.go
	scp ./glitch-image.go root@161.35.36.3:/root/practice/glitch-image.go
	@echo "Files uploaded successfully!"

# SSH into server and navigate to practice folder
enter-server:
	ssh root@161.35.36.3

# Run server on remote (via SSH)
run-server:
	ssh root@161.35.36.3 "cd /root/practice && go run main.go"

# Run client locally
run-client: upload-server
	@for file in tcpserver*; do \
		if [ -f "$$file" ] && [ "$$file" != "tcpserver.txt" ]; then \
			mv "$$file" tcpserver.txt; \
		fi; \
	done
	go run .
	@if [ -f "tcpserver.txt" ]; then mv tcpserver.txt tcpserver.go; fi

# Clean up generated files locally
clean:
	rm -f blurred_image_file.png
	rm -f image_file.png
