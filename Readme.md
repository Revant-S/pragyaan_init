# Golang MVC Backend

## Cloning repo and setting up .env

1. **Clone the Repository**
    ```sh
    git clone https://github.com/Revant-S/pragyaan_init.git
    cd golang-mvc-backend
    ```

2. **Create and Configure `.env` File**
    ```sh
    cp .env.example .env
    ```

    Open the `.env` file and set the appropriate environment variables.

## Local Setup

1. **Run Initial Setup**
    ```sh
    make initial-setup
    ```

2. **Run the Server**
    ```sh
    make run
    ```
## Docker

Start the service

```sh
docker compose up -d
```

