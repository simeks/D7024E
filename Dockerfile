FROM google/golang

WORKDIR /gopath/src/github.com/simeks/D7024E
ADD . /gopath/src/github.com/simeks/D7024E/

RUN go get github.com\nu7hatch\gouuid

RUN go get github.com/simeks/D7024E

CMD []
ENTRYPOINT ["/gopath/bin/github.com/simeks/D7024E"]