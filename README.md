# sniper-go

A web interface for [Sn1per community edition](https://github.com/1N3/Sn1per).

## Features

- Nice UI (thanks to [Bootstrap](https://getbootstrap.com/))
- Scan reports (with download)

## Installation

With docker (recommended):
```bash
docker build -t sniper-go .
```

Or you can install in your system, but you need to install the dependencies:

- [Sn1per](https://github.com/1N3/Sn1per)
- [Go](https://golang.org/)
- [Yarn](https://yarnpkg.com/) (or npm)
- [Make](https://www.gnu.org/software/make/)

After that, you can build the project:

```bash
make setup && make build
```

## Usage

With docker (recommended):
```bash
docker run -p 8080:8080 --name sniper-go sniper-go
```
this will run the container in the background, you can access the web interface at http://localhost:8080.

>Note: If you want to make the data persistent, you can mount a volume to the container in a way similar to this:
>```bash
>docker run -p 8080:8080 --name sniper-go -v /path/to/data:/usr/share/sniper/loot sniper-go
>```

From compiled binary:
```bash
./bin/sniper-go
```
