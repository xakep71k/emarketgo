FROM ubuntu:latest
COPY data.txt /data/products.txt
COPY web_root /www
WORKDIR /sandbox
COPY emarket/emarket .
EXPOSE 8080
CMD ["/sandbox/emarket", "--web-root=/www", "--listen", ":8080", "--data", "/data/products.txt"]
