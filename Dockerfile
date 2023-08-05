FROM xer0dayz/sn1per

LABEL maintainer="matteospanio"

WORKDIR /sniper-go

COPY . .

RUN go mod download

RUN apt update && apt install -y nodejs npm

RUN npm install

RUN npm run build

RUN make build

EXPOSE 8080

ENTRYPOINT ["./bin/sniper-go"]