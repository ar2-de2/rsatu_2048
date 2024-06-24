FROM golang:1.23rc1-bullseye
COPY ./ /opt/rsatu_2048
WORKDIR /opt/rsatu_2048
RUN cd /opt/rsatu_2048 && \
    go mod tidy
CMD ["go", "run", "main.go"]
