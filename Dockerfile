FROM alpine:latest
ADD https://github.com/antham/doc-hunt/releases/download/v2.1.1/doc-hunt_linux_amd64 /usr/bin/doc-hunt
RUN chmod +x /usr/bin/doc-hunt
