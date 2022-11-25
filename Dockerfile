FROM scratch
COPY go-cistercian /
ENTRYPOINT ["/go-cistercian"]
