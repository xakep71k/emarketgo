FROM golang:latest
WORKDIR /build
COPY emarket .
RUN go install -v ./...
EXPOSE 8080
CMD ["emarket", "--web-root=/home/alek/emarketgo/web_root"]
