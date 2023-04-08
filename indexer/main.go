package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/joho/godotenv"
)

var (
	es *elasticsearch.Client
	indexed uint64
) 



type TMBDMovie struct {
	TMDbID string `json:"TMDb_Id"`
	IMDbID string `json:"IMDb_Id"`
	Title  string `json:"Title"`
	OriginalTitle string `json:"Original_Title"`
	Overview string `json:"Overview"`
	Genres []string `json:"Genres"`
	Cast []string `json:"Cast"`
	Crew []string `json:"Crew"`
	Collection string `json:"Collection"`
	ReleaseDate string `json:"Release_Date"`
	ReleaseStatus string `json:"Release_Status"`
	OriginalLanguage string `json:"Original_Language"`
	LanguagesSpoken []string `json:"Languages_Spoken"`
	Runtime string `json:"Runtime"`
	Tagline string `json:"Tagline"`
	Popularity string `json:"Popularity"`
	RatingAverage string `json:"Rating_Average"`
	RatingCount string `json:"Rating_Count"`
	ProductionCompanies []string `json:"Production_Companies"`
	CountryOfOrigin string `json:"Country_Of_Origin"`
	Budget float64 `json:"Budget"`
	Revenue float64 `json:"Revenue"`
}


func OnSuccess(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
	atomic.AddUint64(&indexed, 1)
}

func OnFailure(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
	} else {
		log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
	}
}

func bulkLoad(movies []TMBDMovie) (error) {

	bi, err := esutil.NewBulkIndexer(
		esutil.BulkIndexerConfig{
			Index: "tmdb",
			Client: es,
			NumWorkers: 4,
			FlushBytes: 5e+6,
			FlushInterval: 30 * time.Second,
		},
	)

	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	// Re-create the index
	res, err := es.Indices.Delete([]string{"tmdb"}, es.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res.Body.Close()

	res, err = es.Indices.Create("tmdb")
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()
	
	start := time.Now().UTC()

	context := context.Background()

	for _, movie := range movies {
		data, err := json.Marshal(movie)
		if err != nil {
			log.Fatalf("Error marshalling movie: %s", err)
		}

		err = bi.Add(
			context,
			esutil.BulkIndexerItem{
				Action: "index",
				DocumentID: movie.IMDbID,
				Body:   bytes.NewReader(data),
				OnSuccess: OnSuccess,
				OnFailure: OnFailure,
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	if err := bi.Close(context); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%d] documents with [%d] errors in %s (%d docs/sec)",
			int64(biStats.NumFlushed),
			int64(biStats.NumFailed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%d] documents in %s (%d docs/sec)",
			int64(biStats.NumFlushed),
			dur.Truncate(time.Millisecond),
			int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
		)
	}

	return nil
}

func trimSpaces(line []string) []string {
	for i, v := range line {
		line[i] = strings.TrimSpace(v)
	}
	return line
}

func parseMovie(line []string) (TMBDMovie, error) {
	if len(line) != 22 {
		return TMBDMovie{}, errors.New("Invalid line")
	}

	movie := TMBDMovie{
		TMDbID: line[0],
		IMDbID: line[1],
		Title: line[2],
		OriginalTitle: line[3],
		Overview: line[4],
		Genres: trimSpaces(strings.Split(line[5], "|")),
		Cast: trimSpaces(strings.Split(line[6], "|")),
		Crew: trimSpaces(strings.Split(line[7], "|")),
		Collection: line[8],
		ReleaseDate: line[9],
		ReleaseStatus: line[10],
		OriginalLanguage: line[11],
		LanguagesSpoken: trimSpaces(strings.Split(line[12], "|")),
		Runtime: line[13],
		Tagline: line[14],
		Popularity: line[15],
		RatingAverage: line[16],
		RatingCount: line[17],
		ProductionCompanies: trimSpaces(strings.Split(line[18], "|")),
		CountryOfOrigin: line[19],
	}

	var err error
	movie.Budget, err = strconv.ParseFloat(line[20], 64)
	if err != nil {
		return TMBDMovie{}, errors.New("Invalid budget")
	}

	movie.Revenue, err = strconv.ParseFloat(line[21], 64)
	if err != nil {
		return TMBDMovie{}, errors.New("Invalid revenue")
	}

	return movie, nil
}


func parseMovies(data [][]string) []TMBDMovie {
    movies := make([]TMBDMovie, 0, len(data)-1)

    for i, line := range data {
        if i > 0 { // omit header line
            movie, err := parseMovie(line)
			if err != nil {
				continue
			}
			
			movies = append(movies, movie)
        }


    }
    return movies
}

func main() {
    // open file
    f, err := os.Open("../datasets/TMDB_10000_Popular_Movies.csv")
    if err != nil {
        log.Fatal(err)
    }

    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)
    data, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    // convert records to array of structs
    movies := parseMovies(data)

	err = bulkLoad(movies)
	if err != nil {
		log.Fatalf("Error bulk loading movies: %s", err)
	}
}

func init() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cert, err := os.ReadFile(os.Getenv("ELASTICSEARCH_CA_CERT"))
	if err != nil {
		log.Fatalf("Error reading CA cert: %s", err)
		return
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
		Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
		CACert: cert,
	}

	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
	log.Fatalf("Error creating the client: %s", err)
	}

	// Check if the Elasticsearch cluster is running
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
}