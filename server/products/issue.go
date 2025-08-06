package products

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	utils "github.com/johnietre/utils/go"
)

const (
	issueTimeout = time.Minute * 5
)

var (
	issuesTimes = *utils.NewSyncMap[string, time.Time]()
)

func issueTimeChecker() {
	ticker := time.NewTicker(time.Minute)
	for now := range ticker.C {
		issuesTimes.Range(func(ip string, deadline time.Time) bool {
			if now.After(deadline) {
				issuesTimes.Delete(ip)
			}
			return true
		})
	}
}

type GetIssuesQuery struct {
	SortBy         *string
	SortDesc       *bool
	FilterProduct  *uint64
	FilterReason   *string
	FilterStarted  *bool
	FilterResolved *bool
}

func GetProductIssues(query GetIssuesQuery) (issues []ProductIssue, e error) {
	whereClauses := []string{}
	if query.FilterProduct != nil {
		// TODO: use app ID?
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(`product_id=%d`, *query.FilterProduct),
		)
	}
	if query.FilterReason != nil {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf(`reason=%q`, *query.FilterReason),
		)
	}
	if query.FilterStarted != nil {
		whereClauses = append(
			whereClauses,
			fmt.Sprint(`started_at>0`),
		)
	}
	if query.FilterResolved != nil {
		whereClauses = append(
			whereClauses,
			fmt.Sprint(`resolved_at>0`),
		)
	}
	whereClause := ""
	if len(whereClauses) != 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	sortClause := " ORDER BY id"
	if query.SortBy != nil {
	}
	if query.SortDesc != nil {
		if *query.SortDesc {
			sortClause += " ASC"
		} else {
			sortClause += " DESC"
		}
	}

	rows, err := productsDb.Query(`SELECT * FROM issues` + whereClause + sortClause)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		issue := ProductIssue{}
		err := rows.Scan(
			&issue.Id, &issue.ProductId, &issue.Email,
			&issue.Reason, &issue.Subject, &issue.Description,
			&issue.CreatedAt, &issue.StartedAt, &issue.ResolvedAt,
			&issue.Ip,
		)
		if err != nil {
			if e == nil {
				e = err
			}
		} else {
			issues = append(issues, issue)
		}
	}
	return
}

// AAIError is an error returned by the AddProductIssue function.
type AAIError struct {
	msg string
}

// Error implements the Error method for the error interface.
func (a *AAIError) Error() string {
	return a.msg
}

// IsAAIError returns whether the error is an AAIError
func IsAAIError(err error) bool {
	var errPtr *AAIError
	return errors.As(err, &errPtr)
}

// AddProductIssue adds an issue for the app.
func AddProductIssue(appName string, issue ProductIssue) (ProductIssue, error) {
	if _, err := mail.ParseAddress(issue.Email); err != nil {
		return issue, &AAIError{"bad email"}
	} else if issue.Reason == "" {
		return issue, &AAIError{"missing reason"}
	} else if issue.Subject == "" {
		return issue, &AAIError{"missing subject"}
	} else if issue.Description == "" {
		return issue, &AAIError{"missing description"}
	} else if len(issue.Description) > MaxDescriptionLen {
		return issue, &AAIError{"description too long"}
	}

	// TODO: Do better
	deadline := time.Now().Add(issueTimeout)
	if _, loaded := issuesTimes.LoadOrStore(issue.Ip, deadline); loaded {
		return issue, &AAIError{"wait a bit before submitting"}
	} else if _, loaded = issuesTimes.LoadOrStore(issue.Email, deadline); loaded {
		return issue, &AAIError{"wait a bit before submitting"}
	}

	products, found := productsPageData.Load().Products, false
	for _, app := range products {
		if app.Name == appName {
			issue.ProductId = app.Id
			found = true
			break
		}
	}
	if !found {
		return issue, &AAIError{fmt.Sprintf("no app with name %s", appName)}
	}
	res, err := productsDb.Exec(
		`INSERT INTO issues(
      product_id,email,
      reason,subject,description,
      created_at,started_at,resolved_at,
      ip)
      VALUES (?,?,?,?,?,?,?,?,?)`,
		issue.ProductId, issue.Email,
		issue.Reason, issue.Subject, issue.Description,
		issue.CreatedAt, 0, 0,
		issue.Ip,
	)
	if err != nil {
		return issue, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return issue, err
	}
	issue.Id = uint64(id)
	return issue, nil
}

// EditProductIssue adds an issue for the app.
func EditProductIssue(issue ProductIssue) error {
	if _, err := mail.ParseAddress(issue.Email); err != nil {
		return &AAIError{"bad email"}
	} else if issue.Reason == "" {
		return &AAIError{"missing reason"}
	} else if issue.Subject == "" {
		return &AAIError{"missing subject"}
	} else if issue.Description == "" {
		return &AAIError{"missing description"}
	} else if len(issue.Description) > MaxDescriptionLen {
		return &AAIError{"description too long"}
	}
	_, err := productsDb.Exec(
		`UPDATE issues
    SET product_id=?,email=?,
      reason=?,subject=?,description=?,
      created_at=?,started_at=?,resolved_at=?,
      ip=?
      WHERE id=?`,
		issue.ProductId, issue.Email,
		issue.Reason, issue.Subject, issue.Description,
		issue.CreatedAt, issue.StartedAt, issue.ResolvedAt,
		issue.Ip, issue.Id,
	)
	return err
}

type ProductIssue struct {
	Id          uint64 `json:"id,omitempty"`
	ProductId   uint64 `json:"productId"`
	Email       string `json:"email"`
	Reason      string `json:"reason"`
	Subject     string `json:"subject"`
	Description string `jsom:"description"`
	CreatedAt   int64  `json:"createdAt,omitempty"`
	StartedAt   int64  `json:"startedAt,omitempty"`
	ResolvedAt  int64  `json:"resolvedAt,omitempty"`
	Ip          string `json:"ip,omitempty"`
}
