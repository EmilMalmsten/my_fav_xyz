# Top XYZ

![build](https://github.com/EmilMalmsten/my_top_xyz/actions/workflows/docker-publish.yml/badge.svg?event=push)

Create a toplist for anything!
We all have areas in our lives where we are a bit extra passionate. For you it could be horror movies, hamburger restaurants, cat breeds, or anything else.

Top XYZ is a web app that helps you to keep track of your personal favorites, enables you to rank them, and sharing them with others.

## 🚀 Quick Start

Navigate to [topxyz.net](https://topxyz.net), browse existing toplists or create your own.
Creating your own toplists requires registering with your email and a password.

## 🛠️ Built With

-   Go backend API with CRUD operations for toplists and user account management
-   React frontend
-   JWT authentication
-   PostgreSQL database
-   Docker
-   CI/CD with Github actions

## 🤝 Contributing

### Clone the repo

```
git clone https://github.com/EmilMalmsten/my_top_xyz.git
cd my_top_xyz
```

### Run the project locally in docker

1. Make sure docker is installed by typing: `docker -v` in to your terminal.
   If it's not installed please follow the [docker getting started guide.](https://www.docker.com/get-started/)

2. Read .env.example, re-name it to .env and set up the required environment variables.

3. Build the project with `docker compose build`
4. Start with `docker compose up -d`
5. Run database migrations within the backend container. Replace DB_URL with the actual url used in .env
    ```
    docker exec -it backend sh
    migrate -path migrations -database ${DB_URL} up
    ```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
