package journals

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	//utils "github.com/johnietre/utils/go"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//JournalsData = utils.NewAValue[JournalData](JournalData{})
	journalsDb *sql.DB
)

func InitJournals(journalsDir, dbPath string) error {
	if err := openJournalsDB(dbPath); err != nil {
		return err
	}
	return LoadJournalData()
}

func openJournalsDB(dbPath string) (err error) {
	journalsDb, err = sql.Open("sqlite3", dbPath)

	scriptPath := filepath.Join(filepath.Dir(dbPath), "journals.sql")
	if bytes, err := os.ReadFile(scriptPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("journals DB script (%s) not found", scriptPath)
		} else {
			log.Printf("error reading journals DB script (%s): %v", scriptPath, err)
		}
	} else if _, err := journalsDb.Exec(string(bytes)); err != nil {
		log.Printf("error executing journals DB script (%s): %v", scriptPath, err)
	}
	return
}

func LoadJournalData() error {
	/*
	  journalRows, err := appsDb.Query(`SELECT * FROM journals`)
	  if err != nil {
	    return err
	  }
	  defer journalRows.Close()
	  data := JournalData{}
	  for journalRows.Next() {
	    journal := Journal{}
	    if err := journalRows.Scan(
	      &journal.Id, &journal.Title, &journal.Timestamp, &journal.TzOffset,
	    ); err != nil {
	      return err
	    }
	    data.Journals = append(data.Journals, journal)
	  }

	  catRows, err := appsDb.Query(`SELECT * FROM categories`)
	  if err != nil {
	    return err
	  }
	  defer catRows.Close()
	  for catRows.Next() {
	    cat := ""
	    if err := catRows.Scan(&cat); err != nil {
	      return err
	    }
	    data.Categories = append(data.Categories, cat)
	  }

	  JournalsData.Store(data)
	*/
	return nil
}

func journalNotFound() {
}

type JournalsPageData struct {
	Categories []string
	Journals   []Journal
}

func NewJournalsPageData() JournalsPageData {
	return JournalsPageData{
		Categories: []string{"One", "Two", "Three"},
		Journals: []Journal{
			{Id: 1, Title: "Yes"},
			{Id: 2, Title: "No"},
			{Id: 3, Title: "Maybe"},
		},
	}
}

type Journal struct {
	Id        uint64 `json:"id"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
	TzOffset  int    `json:"tzOffset"`
}

type JournalIndex struct {
	Categories []string
}
