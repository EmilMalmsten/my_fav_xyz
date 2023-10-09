# Top XYZ

Create a toplist for anything!
We all have areas in our lives where we are a bit extra passionate. For you it could be horror movies, hamburger restaurants, cat breeds, or anything else.

Top XYZ is a web app that helps you to keep track of your personal favorites, enables you to rank them, and sharing them with others.

## 🚀 Quick Start

Navigate to [topxyz.net](https://topxyz.net), browse existing toplists or create your own.
Creating your own toplists requires registering with your email and a password.

## 🛠️ Technical features

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

3. Run the project with `docker compose up`

### Run tests

1.  cd in to the backend directory
2.  `go test ./...`

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
