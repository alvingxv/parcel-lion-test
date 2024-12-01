package movierepo

import (
	"context"
	"database/sql"
	"errors"
	"lion-parcel-test/constant"
	"lion-parcel-test/internal/interfaces/adapter"
	"lion-parcel-test/internal/interfaces/repository"
	"lion-parcel-test/pkg/errs"
	"math"

	"go.elastic.co/apm/v2"
)

type movieRepository struct {
	database adapter.DatabaseClient
}

func NewMovieRepository(database adapter.DatabaseClient) repository.MovieRepository {
	return &movieRepository{
		database: database,
	}
}

func (rp *movieRepository) InsertMovieToDB(ctx context.Context, Title string, Description string, Duration int, Artist string, Genre string, FileName string) errs.MessageErr {
	apmSpan, ctx := apm.StartSpan(ctx, "InsertUserToDB", "Repository")
	defer apmSpan.End()

	insertMovieQuery := `INSERT INTO movies (title, description, duration, artists, genres, watch_url) VALUES (?, ?, ?, ?, ?, ?);`

	watchUrl := "localhost:8080/movies/" + FileName

	result := rp.database.Execute(ctx, insertMovieQuery, Title, Description, Duration, Artist, Genre, watchUrl)
	if result.Error != nil {
		return errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			result.Error.Error(),
		)
	}

	return nil
}

func (rp *movieRepository) UpdateMovieToDB(ctx context.Context, Id string, Title string, Description string, Duration int, Artist string, Genre string, FileName string) errs.MessageErr {
	apmSpan, ctx := apm.StartSpan(ctx, "UpdateMovieToDB", "Repository")
	defer apmSpan.End()

	updateMovieQuery := `UPDATE movies SET title = ?, description = ?, duration = ?, artists = ?, genres = ?, watch_url = ? WHERE id = ?;`

	watchUrl := "localhost:8080/movies/" + FileName

	result := rp.database.Execute(ctx, updateMovieQuery, Title, Description, Duration, Artist, Genre, watchUrl, Id)
	if result.Error != nil {
		return errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			result.Error.Error(),
		)
	}

	return nil
}

func (rp *movieRepository) GetMostViewedMovieFromDB(ctx context.Context) (repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetMostViewedMovieFromDB", "Repository")
	defer apmSpan.End()

	getTopMovieQuery := `SELECT id, title, description, duration, artists, genres, watch_url, views_count FROM movies ORDER BY views_count DESC LIMIT 1;`

	var movie repository.Movie

	row := rp.database.QueryRow(ctx, getTopMovieQuery)
	err := row.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.Genre, &movie.WatchUrl, &movie.Views)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.Movie{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return repository.Movie{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}
	return movie, nil
}

func (rp *movieRepository) GetMostViewedGenreFromDB(ctx context.Context) (repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetMostViewedGenreFromDB", "Repository")
	defer apmSpan.End()

	getTopMovieQuery := `SELECT genres, SUM(views_count) AS total_views FROM movies GROUP BY genres ORDER BY total_views DESC LIMIT 1;`

	var movie repository.Movie

	row := rp.database.QueryRow(ctx, getTopMovieQuery)
	err := row.Scan(&movie.Genre, &movie.Views)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.Movie{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return repository.Movie{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}
	return movie, nil
}

func (rp *movieRepository) GetMoviesFromDB(ctx context.Context, page int, pageSize int) ([]repository.Movie, repository.MoviePaginationMetadata, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetMoviesFromDB", "Repository")
	defer apmSpan.End()

	offset := (page - 1) * pageSize

	// Get total number of movies
	var totalItems int
	countQuery := `SELECT COUNT(*) FROM movies`
	row := rp.database.QueryRow(ctx, countQuery)
	err := row.Scan(&totalItems)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.MoviePaginationMetadata{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return nil, repository.MoviePaginationMetadata{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}

	getMoviesQuery := `SELECT id, title, description, duration, artists, genres, watch_url, views_count FROM movies LIMIT ? OFFSET ?`
	movies := make([]repository.Movie, 0)

	rows, err := rp.database.QueryRows(ctx, getMoviesQuery, pageSize, offset)
	if err != nil {
		return nil, repository.MoviePaginationMetadata{}, errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			err.Error(),
		)
	}

	for rows.Next() {
		var movie repository.Movie

		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.Genre, &movie.WatchUrl, &movie.Views)
		if err != nil {
			return nil, repository.MoviePaginationMetadata{}, errs.NewCustomErrs(
				"Failed Scan",
				"FD",
				err.Error(),
			)
		}

		movies = append(movies, movie)

	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return movies, repository.MoviePaginationMetadata{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}, nil

}

func (rp *movieRepository) SearchMoviesFromDB(ctx context.Context, keyword string) ([]repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "SearchMoviesFromDB", "Repository")
	defer apmSpan.End()

	getMoviesQuery := `
    SELECT id, title, description, duration, artists, genres, watch_url, views_count
    FROM movies
    WHERE title LIKE ? OR description LIKE ? OR artists LIKE ? OR genres LIKE ?;`
	movies := make([]repository.Movie, 0)

	searchTerm := "%" + keyword + "%"
	rows, err := rp.database.QueryRows(ctx, getMoviesQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			err.Error(),
		)
	}

	for rows.Next() {
		var movie repository.Movie

		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.Genre, &movie.WatchUrl, &movie.Views)
		if err != nil {
			return nil, errs.NewCustomErrs(
				"Failed Scan",
				"FD",
				err.Error(),
			)
		}

		movies = append(movies, movie)

	}

	return movies, nil
}

func (rp *movieRepository) InsertVoteToDB(ctx context.Context, userId int, movieId int) errs.MessageErr {
	apmSpan, ctx := apm.StartSpan(ctx, "InsertVoteToDB", "Repository")
	defer apmSpan.End()

	insertVoteQuery := `INSERT INTO votes (user_id, movie_id) VALUES (?, ?);`

	result := rp.database.Execute(ctx, insertVoteQuery, userId, movieId)
	if result.Error != nil {

		if result.Error.Error() == constant.DuplicateConstraintError {
			return errs.NewCustomErrs(
				"Already Voted",
				"AV",
				result.Error.Error(),
			)
		}

		return errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			result.Error.Error(),
		)
	}

	return nil
}

func (rp *movieRepository) DeleteVoteFromDB(ctx context.Context, userId int, movieId int) errs.MessageErr {
	apmSpan, ctx := apm.StartSpan(ctx, "InsertVoteToDB", "Repository")
	defer apmSpan.End()

	deleteVoteQuery := `DELETE FROM votes WHERE user_id = ? AND movie_id = ?;`

	result := rp.database.Execute(ctx, deleteVoteQuery, userId, movieId)
	if result.Error != nil {
		return errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			result.Error.Error(),
		)
	}

	return nil
}

func (rp *movieRepository) GetAllVotedMoviesByUserIdFromDb(ctx context.Context, userId int) ([]repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetAllVotedMoviesByUserIdFromDb", "Repository")
	defer apmSpan.End()

	getVotedMoviesQuery := `
    SELECT m.id, m.title, m.description, m.duration, m.artists, m.genres, m.watch_url, m.views_count
	FROM movies m
	JOIN votes v ON m.id = v.movie_id
	WHERE v.user_id = ?;`

	movies := make([]repository.Movie, 0)

	rows, err := rp.database.QueryRows(ctx, getVotedMoviesQuery, userId)
	if err != nil {
		return nil, errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			err.Error(),
		)
	}

	for rows.Next() {
		var movie repository.Movie

		err = rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.Genre, &movie.WatchUrl, &movie.Views)
		if err != nil {
			return nil, errs.NewCustomErrs(
				"Failed Scan",
				"FD",
				err.Error(),
			)
		}

		movies = append(movies, movie)

	}

	return movies, nil
}

func (rp *movieRepository) GetMostVotedMovieFromDB(ctx context.Context) (repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetMostVotedMovieFromDB", "Repository")
	defer apmSpan.End()

	getTopMovieQuery := `
	SELECT m.id, m.title, m.description, m.duration, m.artists, m.genres, m.watch_url, m.views_count, COUNT(v.movie_id) AS vote_count
	FROM movies m
	JOIN votes v ON m.id = v.movie_id
	GROUP BY m.id
	ORDER BY vote_count DESC
	LIMIT 1;
	`

	var movie repository.Movie

	row := rp.database.QueryRow(ctx, getTopMovieQuery)
	err := row.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.Genre, &movie.WatchUrl, &movie.Views, &movie.Vote)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.Movie{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return repository.Movie{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}
	return movie, nil
}

func (rp *movieRepository) GetMostVotedGenreFromDB(ctx context.Context) (repository.Movie, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetMostVotedGenreFromDB", "Repository")
	defer apmSpan.End()

	getTopMovieQuery := `
	SELECT m.genres, COUNT(v.movie_id) AS vote_count
	FROM movies m
	JOIN votes v ON m.id = v.movie_id
	GROUP BY m.genres
	ORDER BY vote_count DESC
	LIMIT 1;
	`

	var movie repository.Movie

	row := rp.database.QueryRow(ctx, getTopMovieQuery)
	err := row.Scan(&movie.Genre, &movie.Vote)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.Movie{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return repository.Movie{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}
	return movie, nil
}
