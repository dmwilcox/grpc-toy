package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DefaultLogger struct {
        logger *log.Logger
        verbose bool
}

func (l DefaultLogger) Print(v ...interface{}) {
        l.logger.Print(v...)
}

func (l DefaultLogger) Printf(format string, v ...interface{}) {
        l.logger.Printf(format, v...)
}

func (l DefaultLogger) Println(v ...interface{}) {
        l.logger.Println(v...)
}

func (l DefaultLogger) Debug(v ...interface{}) {
        if l.verbose {
                l.logger.Print(v...)
        }
}

func (l DefaultLogger) Debugf(format string, v ...interface{}) {
        if l.verbose {
                l.logger.Printf(format, v...)
        }
}

func (l DefaultLogger) Debugln(v ...interface{}) {
        if l.verbose {
                l.logger.Println(v...)
        }
}

func getDefaultLogger(verbose bool) DefaultLogger {
        l := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
        return DefaultLogger{logger: l, verbose: verbose}
}


func main() {
	var verbose bool

	flag.BoolVar(
                &verbose,
                "verbose",
                false,
                "Output debug log messages.",
        )
        flag.Parse()
	logger := getDefaultLogger(verbose)

	user := os.GetEnv("PGUSER")
	password := os.GetEnv("PGPASSWORD")
	database := os.GetEnv("PGDATABASE")
	host := os.GetEnv("PGHOST")
	port := os.GetEnv("PGPORT")

	connectStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	        host,
		port,
		user,
		password,
		database,
	)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}
