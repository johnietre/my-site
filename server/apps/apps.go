package apps

import (
	"database/sql"
	"errors"

	utils "github.com/johnietre/utils/go"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MaxDescriptionLen = 1000
)

var (
	appsPageData = utils.NewAValue[AppsPageData](AppsPageData{})
	appsDb       *sql.DB

	ErrNotFound = errors.New("not found")
)

func InitApps(dbPath string) error {
	if err := openAppsDB(dbPath); err != nil {
		return err
	}
	go issueTimeChecker()
	return LoadAppData()
}

func openAppsDB(dbPath string) (err error) {
	appsDb, err = sql.Open("sqlite3", dbPath)
	return
}

// GetAppById attempts to get an app by its ID.
func GetAppById(id uint64) (App, error) {
	row := appsDb.QueryRow(`SELECT * FROM apps WHERE id=?`, id)
	app := App{}
	err := row.Scan(
		&app.Id, &app.Name, &app.Description,
		&app.Webpage, &app.OnAppStore, &app.OnPlayStore,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return app, nil
}

// GetApps gets all apps.
func GetApps() ([]App, error) {
	const query = `SELECT * FROM apps`
	rows, err := appsDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	apps := []App{}
	for rows.Next() {
		app := App{}
		if err := rows.Scan(
			&app.Id, &app.Name, &app.Description,
			&app.Webpage, &app.OnAppStore, &app.OnPlayStore,
		); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// AddApp adds an app to the database, returning the app with the ID set.
func AddApp(app App) (App, error) {
	res, err := appsDb.Exec(
		`INSERT INTO apps(
      name,description,webage,on_app_store,on_play_store
    ) VALUES (?,?,?,?,?)`,
		app.Name, app.Description, app.Webpage, app.OnAppStore, app.OnPlayStore,
	)
	if err != nil {
		return app, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return app, err
	}
	app.Id = uint64(id)
	return app, nil
}

// EditApp updates the app with the given ID with the passed app.
func EditApp(app App) error {
	_, err := appsDb.Exec(
		`UPDATE apps
    SET name=?,description=?,webpage=?,on_app_store=?,on_play_store=?
    WHERE id=?`,
		app.Name, app.Description,
		app.Webpage, app.OnAppStore, app.OnPlayStore,
		app.Id,
	)
	return err
}

// LoadAppData loads the apps page data.
func LoadAppData() error {
	apps, err := GetApps()
	if err != nil {
		return err
	}
	appsPageData.Store(AppsPageData{
		Apps:              apps,
		MaxDescriptionLen: MaxDescriptionLen,
	})
	return nil
}

type AppsPageData struct {
	Apps              []App
	MaxDescriptionLen int
}

func NewAppsPageData() AppsPageData {
	return appsPageData.Load()
}

type App struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Webpage     string `json:"webpage"`
	OnAppStore  bool   `json:"onAppStore"`
	OnPlayStore bool   `json:"onPlayStore"`
}
