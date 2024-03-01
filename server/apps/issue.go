package apps

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
  issueTimeout = time.Second
)

var (
  issuesTimes = *utils.NewSyncMap[string, time.Time]()
)

func issueTimeChecker() {
  for {
    time.Sleep(time.Minute)
    now := time.Now()
    issuesTimes.Range(func(ip string, deadline time.Time) bool {
      if now.After(deadline) {
        issuesTimes.Delete(ip)
      }
      return true
    })
  }
}

type GetIssuesQuery struct {
  SortBy *string
  SortDesc *bool
  FilterApp *uint64
  FilterReason *string
  FilterRepliedTo *bool
}

func GetAppIssues(query GetIssuesQuery) (issues []AppIssue, e error) {
  whereClauses := []string{}
  if query.FilterApp != nil {
    // TODO: use app ID?
    whereClauses = append(
      whereClauses,
      fmt.Sprintf(`app_id=%d`, *query.FilterApp),
    )
  }
  if query.FilterReason != nil {
    whereClauses = append(
      whereClauses,
      fmt.Sprintf(`reason=%q`, *query.FilterReason),
    )
  }
  if query.FilterRepliedTo != nil {
    val := 0
    if *query.FilterRepliedTo {
      val = 1
    }
    whereClauses = append(
      whereClauses,
      fmt.Sprintf(`replied_to=%d`, val),
    )
  }
  whereClause := ""
  if len(whereClauses) != 0 {
    whereClause = " WHERE" + strings.Join(whereClauses, " AND ")
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

  rows, err := appsDb.Query(`SELECT * FROM issues`+whereClause+sortClause)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, nil
    }
    return nil, err
  }
  defer rows.Close()
  for rows.Next() {
    issue := AppIssue{}
    err := rows.Scan(
      &issue.Id, &issue.AppId, &issue.Email, &issue.Reason,
      &issue.Reason, &issue.Subject, &issue.Description, &issue.RepliedTo,
      &issue.Ip, &issue.Timestamp,
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

// AAIError is an error returned by the AddAppIssue function.
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

// AddAppIssue adds an issue for the app.
func AddAppIssue(appName string, issue AppIssue) (AppIssue, error) {
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

  apps, found := appsPageData.Load().Apps, false
  for _, app := range apps {
    if app.Name == appName {
      issue.AppId = app.Id
      found = true
      break
    }
  }
  if !found {
    return issue, &AAIError{fmt.Sprintf("no app with name %s", appName)}
  }
  res, err := appsDb.Exec(
    `INSERT INTO issues(app_id,email,reason,subject,description,replied_to,ip,timestamp)
      VALUES (?,?,?,?,?,?,?,?)`,
    issue.AppId, issue.Email,
    issue.Reason, issue.Subject, issue.Description,
    false,
    issue.Ip, issue.Timestamp,
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

type AppIssue struct {
  Id uint64
  AppId uint64
  Email string
  Reason string
  Subject string
  Description string
  RepliedTo bool
  Ip string
  Timestamp int64
}
