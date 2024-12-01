# Parcel Lion Technical Test - Alvian Ghifari

This submission was completed in 5 hours. Since I had other things to do, I couldn't create the unit tests. Please note that the submitted movie is not an actual movie; I simulated the uploaded movie by submitting an image. The submission wasn't planned thoroughly, as everything was aimed at simplicity. Also, I completely forgot to commit my work as I went along, so the commits you'll see are only one, haha xd.

### Built With
Golang, Fiber, SQLite (for simplicity sake), Viper.

## Technical Decision
- Used golang (Fiber) and microservices for performance and delivery speed,
- Project bootstrapped with APM for observability, robust logging that can be used to monitoring from Elastic, graceful shutdown, circuit breaker for each external service, clean code with granular testability (small unit test), general interface for multiple implementation option, development and production application configuration, stateless application that can be horizontally scaled using container orchestration.
- Using goroutines when calling multiple other services (if the calls not dependant one to another)

## APIs

### All Users
- POST /api/v1/register — User registration (userHandler.Register)
- POST /api/v1/login — User login (userHandler.Login)
- GET /api/v1/movies — Get all movies (movieHandler.GetMovies)
- GET /api/v1/movies/search — Search movies (movieHandler.SearchMovies)
### Admin (Requires Admin Authentication)
- POST /api/v1/admin/movies — Create a movie (movieHandler.CreateMovie)
- PUT /api/v1/admin/movies/:id — Update a movie (movieHandler.UpdateMovie)
- GET /api/v1/admin/movies/most_viewed — Get most viewed movies (movieHandler.MostViewed)
- GET /api/v1/admin/movies/most_viewed_genre — Get most viewed movies by genre - (movieHandler.MostViewedGenre)
- GET /api/v1/admin/movies/most_voted — Get most voted movies (movieHandler.MostVoted)
- GET /api/v1/admin/movies/most_voted_genre — Get most voted movies by genre (movieHandler.MostVotedGenre)
### Authenticated Users (Requires Authentication)
- POST /api/v1/movies/vote — Vote for a movie (movieHandler.VoteMovie)
- POST /api/v1/movies/unvote — Unvote a movie (movieHandler.UnvoteMovie)
- GET /api/v1/movies/votes — Get voted movies (movieHandler.VotedMovies)

### Postman
[Postman collection](lion-parcel-test-Alvian%20Ghifari.postman_collection.json)





## Database Design.

### users
```sql
CREATE TABLE
  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  )
```

### movies
```sql
CREATE TABLE
  movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL,
    artists TEXT,
    genres TEXT,
    watch_url TEXT,
    views_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  )
```

### votes
```sql
CREATE TABLE
  votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, movie_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE
  )
```

## Project Structure
```
|   go.mod
|   go.sum
|   movies.db -> sqlite db
|   
+---cmd -> commands, entry point of application
|       main.go
|
+---config -> application config
|       config.go
|       dev.yaml
|
+---constant -> all constant
|       app.go
|
+---internal
|   +---adapters -> all related to outside service will be handled here
|   |   +---database
|   |   |   \---sqlite
|   |   |           sqlite.go
|   |   |
|   |   \---micro
|   +---app -> dependency injection stuff and adapter initialization
|   |       dependencies.go
|   |       main.go
|   |       repositories.go
|   |       usecases.go
|   |
|   +---delivery -> delivery method, could be http, grpc, kafka, etc.
|   |   +---grpc
|   |   \---http
|   |           http.go
|   |           movie.go
|   |           user.go
|   |
|   +---interfaces -> all the interfaces will be gathered here
|   |   +---adapter
|   |   |       database.go
|   |   |
|   |   +---delivery
|   |   |       movie.go
|   |   |       user.go
|   |   |
|   |   +---repository
|   |   |       movie.go
|   |   |       user.go
|   |   |
|   |   \---usecase
|   |           movie.go
|   |           user.go
|   |
|   +---repository -> data access layer
|   |   +---movie
|   |   |       movie.go
|   |   |
|   |   \---user
|   |           user.go
|   |
|   \---usecase -> usecases or all the business process
|       +---movie -> movie related usecase
|       |       create_movie.go
|       |       get_movies.go
|       |       most_viewed.go
|       |       most_viewed_genre.go
|       |       most_voted.go
|       |       most_voted_genre.go
|       |       movie.go
|       |       search_movies.go
|       |       unvote_movie.go
|       |       update_movie.go
|       |       voted_movies.go
|       |       vote_movie.go
|       |
|       \---user -> user related usecase
|               login.go
|               populate_session.go
|               register.go
|               user.go
|
+---logs -> application logs
|
+---movies -> uploaded movies
|       logo-bundar.png
|
\---pkg -> all the global package
    +---dto
    |       dto.go
    |
    +---errs
    |       errs.go
    |
    +---httpclient
    |       httpclient.go
    |
    +---log
    |       logger.go
    |
    \---middleware
            setup.go


```