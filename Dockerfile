FROM xer0dayz/sn1per

LABEL maintainer="matteospanio"

WORKDIR /sniper-go

COPY . .

RUN apt update && apt install -y nodejs npm && npm install -g yarn

RUN make setup && make build

RUN mkdir "/usr/share/sniper/loot/workspace"

VOLUME [ "/usr/share/sniper/loot" ]

EXPOSE 8080

CMD ["./bin/sniper-go", "-port", "8080", "-mode", "debug"]