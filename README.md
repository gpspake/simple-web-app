# Simple Web App

A web application built with Go, Echo, SQLite, and HTMX. This app demonstrates basic routing, pagination, and full text search.

Demo: https://simple-web-app.georgespake.com

---

## Features

- **Releases Management**: View a paginated, searchable list of music releases with details like release year and associated artists.
- **Full Text Search**: Uses the FTS5 sqlite extension with trigram tokenization.
- **Templating**: Uses Go's `html/template` package for rendering HTML pages.
- **In-Memory Testing**: Comprehensive test coverage with an in-memory SQLite database.

---

## Tech Stack

- **Backend**: [Go](https://golang.org/)
- **Web Framework**: [Echo](https://echo.labstack.com/)
- **Database**: [SQLite](https://sqlite.org/index.html)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Styling**: [Tailwind](https://github.com/golang-migrate/migrate)
- **Partial UI Re-rendering**: [HTMX](https://github.com/golang-migrate/migrate)
- **Testing**: Built-in Go testing framework with [Testify](https://github.com/stretchr/testify)

---

## Prerequisites

- Go 1.20 or later
- SQLite installed locally (for development)
- Docker (optional, for containerized deployment)

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/gpspake/simple-web-app
cd simple-web-app
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run the App
```bash
go run -tags "sqlite_fts5" .
```

### 4. Visit the site

Navigate to [localhost:8086](http://localhost:8086) in a web browser.

---

## Local development

### With Air

Use [Air](https://github.com/air-verse/air) to auto reload when files are changed. 

#### 1. Install Air
Make sure `$GOPATH/bin` is in your Path
```bash
go install github.com/air-verse/air@latest
```

#### 2. Run the app using Air
```bash
air -c .air.toml
```

### With Docker
You can also run the app with docker.
```bash
docker build -t simple-web-app .
docker run -p 8086:8086 simple-web-app
```

Or you can run
```bash
./start.sh
```

---

## Deployment
This app is configured to deploy to a Digital Ocean using github actions. 

_Note: This section provides a general overview but it's not a comprehensive deployment guide. Some additional server configuration may be required._

1. On the Droplet: Create a user for deployment and grant the user access to docker.
    ```
   sudo adduser deploy
   sudo usermod -aG docker deploy
   ```
1. Switch to the user and add the public key to the user's authorized keys
    ```
   sudo su - deploy
   mkdir -p ~/.ssh
   chmod 700 ~/.ssh
   touch ~/.ssh/authorized_keys
   chmod 600 ~/.ssh/authorized_keys
   echo "<YOUR_PUBLIC_KEY_CONTENT>" >> ~/.ssh/authorized_keys
   ```
1. On Github: Go to your repository's Settings > Secrets and variables > Actions and add the following secrets:
    ```
   DIGITALOCEAN_IP: The IP address of your droplet.
    SSH_USERNAME: Your SSH username on the droplet.
    SSH_PRIVATE_KEY: Your private SSH key for authentication.
   ```
1. Run the action. The [deploy workflow](./.github/workflows/deploy.yml)  is configured to run on push to main.
1. Verify the container is running:
    ```
   docker ps
    CONTAINER ID   IMAGE                   COMMAND   CREATED         STATUS         PORTS                                                   NAMES
    19282e9c5697   simple-web-app:latest   "./app"   3 minutes ago   Up 3 minutes   8080/tcp, 0.0.0.0:8080->8086/tcp, [::]:8080->8086/tcp   simple-web-app
   ```
1. Serve the app. Example using Caddy:
    ```
   simple-web-app.georgespake.com {
        reverse_proxy 127.0.0.1:8080

        tls {
                dns digitalocean {env.DO_AUTH_TOKEN}
        }

        log {
                output file /var/log/caddy/access.log
                format json
        }
    }
   ```

---

## Other Stuff

To learn more about how this project was built, check out the commits. Each change is committed with a descriptive message and [gitmoji](https://gitmoji.dev/)

The purpose of this project was to learn more about the stack so some of it was new to me. If I got something wrong, let me know or submit a PR.

I might abandon this but I hope to keep it updated as a starter kit for other projects. 

