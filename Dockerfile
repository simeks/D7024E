FROM google/golang

WORKDIR /gopath/src/github.com/simeks/d7024e
ADD . /gopath/src/github.com/simeks/d7024e/

RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/liamzebedee/go-qrp

RUN go get github.com/simeks/d7024e

CMD []
ENTRYPOINT ["/gopath/bin/d7024e"]