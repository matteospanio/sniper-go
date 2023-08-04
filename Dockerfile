FROM xer0dayz/sn1per

LABEL maintainer="matteospanio"

WORKDIR /sniper-go

COPY . .

RUN go mod download

RUN make build

ENTRYPOINT ["./bin/sniper-go"]