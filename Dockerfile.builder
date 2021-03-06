FROM fedora:32

RUN dnf update -y

RUN dnf install make git awscli gcc -y

RUN curl https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz -o /opt/go1.14.3.linux-amd64.tar.gz && \
tar -xf /opt/go1.14.3.linux-amd64.tar.gz -C /opt/ && \
ln -s /opt/go/bin/go /usr/bin/go && \
ln -s /opt/go/bin/gofmt /usr/bin/gofmt

ENV PATH /go/bin:$PATH
ENV GOPATH /go/
ENV GOCACHE /tmp/



CMD /bin/bash