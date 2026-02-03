package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbMid *sql.DB

func InitDBMaster() {
	var err error

	pass, err := Decrypt(DatabasePass)
	if err != nil {
		Info("ERROR DECRYPT")
		log.Fatal(err.Error())
	}

	dbMid, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DatabaseUser, pass, DatabaseIp, DatabasePort, DatabaseName))
	// dbMid, err = sql.Open("mysql", "optima2:5UjHPwhokKw2dzfq@tcp(149.129.215.223:3306)/DB_MID_GO")
	if err != nil {
		log.Fatal(err)
	}
	if err := dbMid.Ping(); err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(apiKey string) (dbConn *sql.DB, err error) {
	var apiUser string
	var confUser string
	var confPass string
	var confPort string
	var confHost string
	var confDefault string

	err = dbMid.QueryRow("CALL GET_USER_API_DETAIL(?, 'DB')", apiKey).Scan(&apiUser, &confUser, &confPass, &confPort, &confHost, &confDefault)
	if err != nil {
		log.Printf("Failed to execute SP GET_USER_API_DETAIL: %s", err.Error())
		return nil, err
	}

	confPass, err = Decrypt(confPass)
	if err != nil {
		log.Printf("Failed to decrypt conf pass: %s", err.Error())
		return nil, err
	}

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", confUser, confPass, confHost, confPort, confDefault))
	if err != nil {
		return nil, err
	}

	if err := dbConn.Ping(); err != nil {
		return nil, err
	}

	return dbConn, err
}

func GetToken(contentType, apiKey, authorization string) (response string, err error) {
	Info("Execute SP: CALL GET_TOKEN('%s', '%s', '%s')", contentType, apiKey, authorization)

	err = dbMid.QueryRow("CALL GET_TOKEN(?, ?, ?)", contentType, apiKey, authorization).Scan(&response)
	if err != nil {
		log.Printf("Failed to execute stored procedure GET_TOKEN or scan result: %v", err)
		return response, err
	}

	Info("Res Execute SP: %s", response)

	return response, err
}

func GetSecret(apiKey string) (response string, err error) {
	Info("SELECT md5(secret) as secret from users_api WHERE apikey = %s", apiKey)

	err = dbMid.QueryRow("SELECT md5(secret) as secret from users_api WHERE apikey = ?", apiKey).Scan(&response)
	if err != nil {
		log.Printf("Failed to execute query SELECT md5(secret) as secret from users_api WHERE apikey = ? : %v", err)
		return response, err
	}

	Info("Res get secret: %s***", response[:10])

	return response, err
}

func CallSp(apiKey, spName string, spParams []any) (response []map[string]any, err error) {
	dbConn, err := ConnectDB(apiKey)
	if err != nil {
		log.Printf("Failed to connect DB: %s", err.Error())
		return response, err
	}

	sql := fmt.Sprintf("CALL %s(", spName)

	for i := range spParams {
		if i+1 == len(spParams) {
			sql += "?)"
		} else {
			sql += "?, "
		}
	}

	if len(spParams) == 0 {
		sql += ")"
	}

	spParamsJson, err := json.Marshal(spParams)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	Info("Sql: CALL %s(%v)", spName, string(spParamsJson[1:len(spParamsJson)-1]))

	rows, err := dbConn.Query(sql, spParams...)
	if err != nil {
		log.Printf("Failed to execute SP %s: %v", sql, err)
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colCount := len(columns)

	for rows.Next() {
		// Create a slice of any's to hold the values for this row.
		// We use any because we don't know the types.
		values := make([]any, colCount)

		// Create a slice of pointers to the values.
		// This is necessary for rows.Scan().
		scanArgs := make([]any, colCount)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan the row into the pointers
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Create a map[string]any for this row
		rowMap := make(map[string]any)
		for i, colName := range columns {
			// We need to handle NULL values, which scan as nil.
			// And bytes, which we often want as strings.
			val := values[i]

			if b, ok := val.([]byte); ok {
				// If it's bytes, convert to string
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		// Add the row map to our results
		response = append(response, rowMap)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	Info("Res Execute %s: %+v", spName, string(jsonData))

	return response, err
}

func ExecuteQuery(apiKey, sql string, params []any) (response []map[string]any, err error) {
	dbConn, err := ConnectDB(apiKey)
	if err != nil {
		log.Printf("Failed to connect DB: %s", err.Error())
		return response, err
	}

	paramsJson, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	Info("Sql: %s | Params: %+v", sql, string(paramsJson))

	rows, err := dbConn.Query(sql, params...)
	if err != nil {
		log.Printf("Failed to execute SP %s: %v", sql, err)
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colCount := len(columns)

	for rows.Next() {
		// Create a slice of any's to hold the values for this row.
		// We use any because we don't know the types.
		values := make([]any, colCount)

		// Create a slice of pointers to the values.
		// This is necessary for rows.Scan().
		scanArgs := make([]any, colCount)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan the row into the pointers
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Create a map[string]any for this row
		rowMap := make(map[string]any)
		for i, colName := range columns {
			// We need to handle NULL values, which scan as nil.
			// And bytes, which we often want as strings.
			val := values[i]

			if b, ok := val.([]byte); ok {
				// If it's bytes, convert to string
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}

		// Add the row map to our results
		response = append(response, rowMap)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	Info("Res Execute %s: %+v", sql, string(jsonData))

	return response, err
}
