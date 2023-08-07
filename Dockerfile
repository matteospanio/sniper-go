FROM xer0dayz/sn1per

LABEL maintainer="matteospanio"

WORKDIR /sniper-go

COPY . .

RUN apt update && apt install -y nodejs npm && npm install -g yarn

RUN yarn install

RUN yarn build

RUN make build

EXPOSE 8080

CMD ["./bin/sniper-go"]