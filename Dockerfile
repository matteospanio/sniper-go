FROM xer0dayz/sn1per

LABEL maintainer="matteospanio"

WORKDIR /sniper-go

COPY . .

RUN go mod download

RUN make build

EXPOSE 8080

ENTRYPOINT ["./bin/sniper-go"]