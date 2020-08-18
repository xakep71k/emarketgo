FROM golang:latest as modules

ADD emarket/go.mod emarket/go.sum /m/
RUN cd /m && go mod download

FROM golang:latest as builder

COPY --from=modules /go/pkg /go/pkg
RUN useradd -u 10001 emarket

RUN mkdir -p /emarket
ADD emarket /emarket
WORKDIR /emarket

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/emarket

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER emarket

COPY --from=builder /emarket/bin/emarket /emarket
COPY web_root /www
COPY data.txt /data/products.json
EXPOSE 8080

CMD ["/emarket", "--web-root=/www", "--listen", ":8080", "--data", "/data/products.json"]
