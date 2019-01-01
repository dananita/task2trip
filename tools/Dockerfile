FROM alpine:latest

# add certificates for https connections
RUN apk --no-cache add ca-certificates

EXPOSE 8000

# copy
COPY task2trip /bin
