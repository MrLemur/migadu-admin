# Migadu Admin Panel

`migadu-admin` is an application written in Go and Svelte to manage mailboxes for domains hosted with [Migadu](https://www.migadu.com) using the official API.

It currently can add/delete mailboxes and identities for multiple domains.

> It is currently in very early stages and should not be relied on yet.

## Setup

The following environment variables are required before running the Go server. These can be set in the environment or a `.env` file placed in the same directory as the Go server.

| Variable           | Description                               | Example                                          |
| ------------------ | ----------------------------------------- | ------------------------------------------------ |
| MIGADU_API_KEY     | API key for Migadu                        | `378ry743rh387ryhjf78erhgu8ih785e4gedrgerbn78hg` |
| MIGADU_ADMIN_EMAIL | Admin email for Migadu account            | `admin@example.com`                              |
| MIGADU_DOMAINS     | Comma separated list of domains to manage | `example.com,test.com,example.net`               |

1. Change into the `frontend` directory and build the app with `yarn build`
2. In the root of the repo, run `go build`
3. In a new directory, copy `migadu-admin` and the `frontend/build` directory
4. Rename `build` to `frontend`
5. Run `./migadu-admin`
6. If there are no errors, the server is available on `https://localhost:5000`

### Docker

A Dockerfile is also provided for your convenience.

1. In the root of the repo, build the Docker image: `docker build -t miagdu-admin:latest .`
2. Run the image: `docker run -p 5000:5000 --env-file .env docker.io/library/miagdu-admin:latest`
