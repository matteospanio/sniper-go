# sniper-go

A web interface for [Sn1per community edition](https://github.com/1N3/Sn1per).

## Features

- [X] Nice UI (thanks to [Bootstrap](https://getbootstrap.com/))
- [X] Scan reports (with download)
- [ ] Scan history (TODO)
- [ ] Scan scheduling (TODO)
- [ ] User management (TODO)

## Installation

Download this repository and change directory:
```bash
git clone https://github.com/matteospanio/sniper-go.git
cd sniper-go
```

You can install this software in two ways:

- With docker (recommended):
    ```bash
    docker build -t sniper-go .
    ```

- Or you can install in your system from source:

    ```bash
    make install # install the dependencies and build the project
    ```

    > Note: `make install` requires root privileges, to have a more fine-grained control in the installation process, see `make help`.

## Usage

- With docker (recommended):
    ```bash
    docker run -p 8080:8080 --name sniper-go sniper-go
    ```
    this will run the container in the background, you can access the web interface at http://localhost:8080.

    >Note: If you want to make the data persistent, you can mount a volume to the container in a way similar to this:
    >```bash
    >docker run -p 8080:8080 --name sniper-go -v /path/to/data:/usr/share/sniper/loot sniper-go
    >```

- From compiled binary:
    ```bash
    ./bin/sniper-go
    ```
    this process needs to be run as root, you can access the web interface at http://localhost:8080.

## Dependencies

`sniper-go` depends on the following software:

- [Sn1per](https://github.com/1N3/Sn1per)
- [Go](https://golang.org/)
- [Yarn](https://yarnpkg.com/) (or npm)
- [Make](https://www.gnu.org/software/make/)

## Credits

- The project is based on [Sn1per community edition](https://github.com/1N3/Sn1per)
- The UI is based on [Bootstrap](https://getbootstrap.com/), [stimulus.js](https://stimulus.hotwired.dev/) and [Font Awesome](https://fontawesome.com/), and it is built with [webpack](https://webpack.js.org/).
- The backend is written in [Go](https://golang.org/) and it uses [Gin](https://github.com/gin-gonic/gin).
