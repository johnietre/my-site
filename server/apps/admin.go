package apps

// AdminAppsListPageData is the page data for admin/apps/list.
type AdminAppsListPageData struct {
  Apps []App
}

// NewAALPageData queries the database and returns a new page data for
// admin/apps/list.
func NewAALPageData() (data AdminAppsListPageData, err error) {
  apps, err := GetApps()
  if err == nil {
    return data, err
  }
  return AdminAppsListPageData{Apps: apps}, nil
}

// AdminAppsListEditPageData is the page data for admin/apps/list/{id}.
type AdminAppsListEditPageData struct {
  App App
}

// NewAALEPageData queries the database and returns a new page data for
// admin/apps/list/{id}. If the ID is 0, the data for a new app is created.
func NewAALEPageData(id uint64) (data AdminAppsListEditPageData, err error) {
  var app App
  if id == 0 {
    data.App = App{Name: "New App"}
  } else {
    app, err = GetAppById(id)
  }
  return AdminAppsListEditPageData{App: app}, nil
}

// AdminAppsReviewPageData queries the database and returns
type AdminAppsReviewPageData struct {
  Issues []AppIssue
}

func NewAARLPageData() AdminAppsReviewPageData {
  return AdminAppsReviewPageData{
    // TODO
  }
}

type AdminAppsReviewReplyPageData struct {
  Issue AppIssue
  Name string
}

func NewAARRPageData() AdminAppsReviewReplyPageData {
  return AdminAppsReviewReplyPageData{
    // TODO
  }
}
