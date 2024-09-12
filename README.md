
# Kitten Bommm API

This is a RESTful API built using Go and the [Fiber](https://gofiber.io/) web framework. The server integrates with Redis for caching and session management, and the deployment is handled via [Railway](https://railway.app/). Redis is hosted on [Aiven](https://aiven.io/), providing a managed, scalable, and secure solution.

## Overview

- **Framework:** [Fiber](https://gofiber.io/) (Go)
- **Database/Cache:** Redis hosted on [Aiven](https://aiven.io/)
- **Deployment:** [Railway](https://railway.app/)
  
The server is responsible for managing users, games, and game-related data such as leaderboards and moves. 

## [Deployment Link](emitrrkittenboomserver-production.up.railway.app)

### Run Local

- [**Download GO**](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04)
- git clone https://github.com/shantanu200/Emitrr_Kitten_Boom_Server.git
- cd Emitrr_Kitten_Boom_Server
- go mod tidy 
- go run main.go


### Key Features:
- **User Management:** Registration, login, and user detail management via JWT-based authentication.
- **Game Management:** Start new games, store moves, retrieve game history, and fetch leaderboard data.
- **Redis Integration:** Redis is used to efficiently store and retrieve session data, caching, and other temporary information to enhance performance.


## Endpoints

### 1. **Check if Username Exists**

- **URL:** `/username/:username`
- **Method:** `GET`
- **Description:** Checks if a username exists in the system.
- **Parameters:**
  - `username` (string): The username to check.
- **Response:**
  - 200: Username exists.
  - 404: Username not found.

### 2. **Register a User**

- **URL:** `/register`
- **Method:** `POST`
- **Description:** Registers a new user.
- **Request Body:**
  - `username` (string): The desired username.
  - `password` (string): The user's password.
- **Response:**
  - 201: User successfully created.
  - 400: Validation error.

### 3. **Login a User**

- **URL:** `/login`
- **Method:** `POST`
- **Description:** Logs in a user and returns a JWT token.
- **Request Body:**
  - `username` (string): The username.
  - `password` (string): The user's password.
- **Response:**
  - 200: Success and returns a JWT token.
  - 401: Unauthorized, invalid credentials.

### 4. **Get User Details**

- **URL:** `/details`
- **Method:** `GET`
- **Description:** Retrieves the details of the authenticated user.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Response:**
  - 200: Returns user details.
  - 401: Unauthorized.

### 5. **Start a New Game**

- **URL:** `/start`
- **Method:** `POST`
- **Description:** Starts a new game for the authenticated user.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Request Body:** (optional) Game setup details.
- **Response:**
  - 201: Game successfully created.
  - 400: Invalid game setup.

### 6. **Store Game Moves**

- **URL:** `/status/:id`
- **Method:** `PATCH`
- **Description:** Updates the status of a game by storing the player's moves.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Parameters:**
  - `id` (string): The game ID.
- **Request Body:** 
  - `moves` (array): List of game moves.
- **Response:**
  - 200: Game moves successfully updated.
  - 404: Game not found.

### 7. **Get User's Games**

- **URL:** `/userGames`
- **Method:** `GET`
- **Description:** Retrieves all games played by the authenticated user.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Response:**
  - 200: List of games.
  - 401: Unauthorized.

### 8. **Get Game by ID**

- **URL:** `/game/:id`
- **Method:** `GET`
- **Description:** Retrieves the details of a specific game by its ID.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Parameters:**
  - `id` (string): The game ID.
- **Response:**
  - 200: Game details.
  - 404: Game not found.

### 9. **Get Leaderboard**

- **URL:** `/leaderboard`
- **Method:** `GET`
- **Description:** Retrieves the leaderboard for the game.
- **Headers:**
  - `Authorization: Bearer <JWT Token>`
- **Response:**
  - 200: Leaderboard data.

---
