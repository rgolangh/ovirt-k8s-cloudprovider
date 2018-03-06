FROM centos:7

LABEL maintainer="rgolan@redhat.com"

RUN yum install -y epel-release
RUN yum install -y golang git glide ansible && yum clean all

RUN mkdir -p /go/{src,bin,pkg}
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ARG PREFIX=/opt

WORKDIR /go/src/github.com/rgolangh/ovirt-cloudprovider
COPY cmd cmd/
COPY internal internal/
COPY Makefile .
COPY glide.lock .
COPY glide.yaml .

RUN make deps build PREFIX=$PREFIX

ENTRYPOINT ["/opt/ovirt-cloudprovider"]
