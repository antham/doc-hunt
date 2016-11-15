FROM frolvlad/alpine-glibc
RUN apk update && apk add pcre
COPY build/doc-hunt_linux_amd64	/usr/bin/doc-hunt
RUN chmod +x /usr/bin/doc-hunt
WORKDIR /app
ENTRYPOINT ["doc-hunt"]
CMD ["help"]
