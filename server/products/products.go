package products

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	utils "github.com/johnietre/utils/go"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MaxDescriptionLen = 1000
)

var (
	productsPageData = utils.NewAValue(ProductsPageData{})
	productsDb       *sql.DB

	ErrNotFound = errors.New("not found")
)

func InitProducts(dbPath string) error {
	if err := openProductsDB(dbPath); err != nil {
		return err
	}
	go issueTimeChecker()
	go func() {
		t := time.NewTicker(time.Minute)
		for range t.C {
			if err := LoadProductData(); err != nil {
				log.Print("error loading product data: ", err)
			}
		}
	}()
	return LoadProductData()
}

func openProductsDB(dbPath string) (err error) {
	productsDb, err = sql.Open("sqlite3", dbPath)

	scriptPath := filepath.Join(filepath.Dir(dbPath), "products.sql")
	if bytes, err := os.ReadFile(scriptPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("products DB script (%s) not found", scriptPath)
		} else {
			log.Printf("error reading products DB script (%s): %v", scriptPath, err)
		}
	} else if _, err := productsDb.Exec(string(bytes)); err != nil {
		log.Printf("error executing products DB script (%s): %v", scriptPath, err)
	}
	return
}

// GetProductById attempts to get an app by its ID.
func GetProductById(id uint64) (Product, error) {
	row := productsDb.QueryRow(`SELECT * FROM products WHERE id=?`, id)
	app := Product{}
	err := row.Scan(
		&app.Id, &app.Name, &app.Description,
		&app.Webpage, &app.AppStoreLink, &app.PlayStoreLink, &app.Hidden,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return app, nil
}

// GetProducts gets all products.
func GetProducts() ([]Product, error) {
	const query = `SELECT * FROM products`
	rows, err := productsDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []Product{}
	for rows.Next() {
		app := Product{}
		if err := rows.Scan(
			&app.Id, &app.Name, &app.Description,
			&app.Webpage, &app.AppStoreLink, &app.PlayStoreLink,
			&app.Hidden,
		); err != nil {
			return nil, err
		}
		products = append(products, app)
	}
	return products, nil
}

// AddProduct adds an app to the database, returning the app with the ID set.
func AddProduct(app Product) (Product, error) {
	res, err := productsDb.Exec(
		`INSERT INTO products(
      name,description,webpage,app_store_link,play_store_link,hidden
    ) VALUES (?,?,?,?,?,?)`,
		app.Name, app.Description, app.Webpage,
		app.AppStoreLink, app.PlayStoreLink,
		app.Hidden,
	)
	if err != nil {
		return app, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return app, err
	}
	app.Id = uint64(id)
	return app, nil
}

// EditProduct updates the app with the given ID with the passed app.
func EditProduct(app Product) error {
	_, err := productsDb.Exec(
		`UPDATE products
    SET name=?,description=?,webpage=?,app_store_link=?,play_store_link=?,hidden=?
    WHERE id=?`,
		app.Name, app.Description,
		app.Webpage, app.AppStoreLink, app.PlayStoreLink,
		app.Hidden, app.Id,
	)
	return err
}

// LoadProductData loads the products page data.
func LoadProductData() error {
	products, err := GetProducts()
	if err != nil {
		return err
	}
	// TODO: something else?
	//utils.FilterSlice(products, func(prod Product) bool { return !prod.Hidden })
	productsPageData.Store(ProductsPageData{
		Products:          products,
		MaxDescriptionLen: MaxDescriptionLen,
	})
	return nil
}

type ProductsPageData struct {
	Products          []Product
	MaxDescriptionLen int
}

func NewProductsPageData() ProductsPageData {
	return productsPageData.Load()
}

type Product struct {
	Id            uint64 `json:"id,omitempty"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Webpage       string `json:"webpage"`
	AppStoreLink  string `json:"appStoreLink,omitempty"`
	PlayStoreLink string `json:"playStoreLink,omitempty"`
	// TODO
	Images []string `json:"images,omitempty"`
	Hidden bool     `json:"hidden,omitempty"`
}
