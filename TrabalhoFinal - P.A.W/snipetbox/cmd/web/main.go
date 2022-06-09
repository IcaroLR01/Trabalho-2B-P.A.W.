package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rmcs87/cc5m/pkg/models/mysql"
)

type application struct{
  errorLog *log.Logger
  infoLog *log.Logger
  snippets *mysql.SnippetModel
}

//curl -i -X POST http://localhost:4000/snippet/create
func main() {
	// nome da flag, valor padra e descrição
	addr := flag.String("addr", ":4000", "Porta da Rede")
  
  dsn := flag.String("dsn",
                     "mz2wv2OwQt:R6BKu1xQCp@tcp(remotemysql.com)/mz2wv2OwQt?parseTime=true", 
                     "MySql DSN")
  
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERRO:\t", log.Ldate|log.Ltime|log.Lshortfile)

  db, err := openDB(*dsn)
  if err != nil{
    errorLog.Fatal(err)
  }
  defer db.Close()
  
  app := &application{
    errorLog: errorLog,
    infoLog: infoLog,
    snippets: &mysql.SnippetModel{DB:db},
  }

  srv := &http.Server{
    Addr: *addr,
    ErrorLog: errorLog,
    Handler: app.routes(),
  }
  
	infoLog.Printf("Inicializando o servidor na porta %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
  db, err := sql.Open("mysql", dsn)
  if err!= nil{
    return nil, err
  }
  if err = db.Ping(); err != nil{
    return nil, err
  }
  return db, nil
}