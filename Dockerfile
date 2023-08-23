FROM jac18281828/godev:latest

ARG PROJECT=ffsid
WORKDIR /workspaces/${PROJECT}
ENV GOMAXPROCS=10
COPY . .
RUN chown -R jac:jac .
USER jac
ENV GOPATH=/workspaces/${PROJECT}

RUN go install -v github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /workspaces/${PROJECT}/src/zkid
RUN go test -v

