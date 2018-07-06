FROM golang:latest 
ENV GOPATH /go
ADD . /go/src/github.com/high-value-team/workshop-kubernetes-setup/workshop-ci
WORKDIR /go/src/github.com/high-value-team/workshop-kubernetes-setup/workshop-ci
RUN go build
RUN mv workshop-ci /workshop-ci
RUN rm -rf vendor
CMD ["/workshop-ci"]