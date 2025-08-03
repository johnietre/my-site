package products

// AdminProductsListPageData is the page data for admin/apps/list.
type AdminProductsListPageData struct {
	Products []Product
}

// NewAALPageData queries the database and returns a new page data for
// admin/apps/list.
func NewAALPageData() (data AdminProductsListPageData, err error) {
	apps, err := GetProducts()
	if err == nil {
		return data, err
	}
	return AdminProductsListPageData{Products: apps}, nil
}

// AdminProductsListEditPageData is the page data for admin/apps/list/{id}.
type AdminProductsListEditPageData struct {
	Product Product
}

// NewAALEPageData queries the database and returns a new page data for
// admin/apps/list/{id}. If the ID is 0, the data for a new app is created.
func NewAALEPageData(id uint64) (data AdminProductsListEditPageData, err error) {
	var app Product
	if id == 0 {
		data.Product = Product{Name: "New Product"}
	} else {
		app, err = GetProductById(id)
	}
	return AdminProductsListEditPageData{Product: app}, nil
}

// AdminProductsReviewPageData queries the database and returns
type AdminProductsReviewPageData struct {
	Issues []ProductIssue
}

func NewAARLPageData() AdminProductsReviewPageData {
	return AdminProductsReviewPageData{
		// TODO
	}
}

type AdminProductsReviewReplyPageData struct {
	Issue ProductIssue
	Name  string
}

func NewAARRPageData() AdminProductsReviewReplyPageData {
	return AdminProductsReviewReplyPageData{
		// TODO
	}
}
