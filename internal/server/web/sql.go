package web

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// hSQL handle sql to json over http request
func (o *Server) hSQL(w http.ResponseWriter, r *http.Request) {
	var err error
	var rows pgx.Rows
	var x []interface{}

	if args, ok := r.URL.Query()["request"]; ok && len(args) == 1 && len(args[0]) > 1 {
		if _, err = io.WriteString(w, fmt.Sprintf("{\n  \"start\":\"%s\",\n  \"request\":\"%s\",", time.Now().Format("2006-01-02 15:04:05"), args[0])); err != nil {
			return
		}
		rows, err = o.dbPool.Query(args[0])
		defer func() {
			if rows != nil {
				rows.Close()
			}
		}()
		if err != nil {
			if _, err = io.WriteString(w, fmt.Sprintf(" \n  \"ok\":\"false\",\n  \"error\":\"%s\",\n  \"result\":\n  [", err.Error())); err != nil {
				return
			}
		} else {
			if _, err = io.WriteString(w, "\n  \"ok\":\"true\",\n  \"result\":\n  ["); err != nil {
				return
			}
			rf := false
			for rows.Next() {
				ff := false
				if rf {
					if _, err = io.WriteString(w, ","); err != nil {
						return
					}
				}
				rf = true
				io.WriteString(w, "\n    {")
				if x, err = rows.Values(); err == nil {
					for i, v := range x {
						if ff {
							if _, err = io.WriteString(w, ","); err != nil {
								return
							}
						}
						ff = true
						if _, err = io.WriteString(w, fmt.Sprintf("\"%v\":\"%v\"", rows.FieldDescriptions()[i].Name, v)); err != nil {
							return
						}
					}
				}
				if _, err = io.WriteString(w, "}"); err != nil {
					return
				}
			}
		}
		io.WriteString(w, fmt.Sprintf("\n  ],\n  \"end\":\"%s\"\n}\n", time.Now().Format("2006-01-02 15:04:05")))
	} else {
		io.WriteString(w, "Uses: sql.html?request=<T-SQL request text>\n")
	}
}
