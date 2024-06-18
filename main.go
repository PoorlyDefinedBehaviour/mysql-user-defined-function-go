package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <mysql.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//export temporal_init
func temporal_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) {
	// Initialization logic if any
}

//export temporal_deinit
func temporal_deinit(initid *C.UDF_INIT) {
	// Cleanup logic if any
}

//export temporal
func temporal(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, is_null *C.char, error *C.char) *C.char {
	ret := "hello world"
	result = C.CString(ret)
	*length = uint64(utf8.RuneCountInString(ret))
	return result
	// Check the number of arguments
	if args.arg_count != 1 {
		C.strcpy(error, C.CString("temporal() requires one argument"))
		return nil
	}

	// Get the input query
	query := C.GoStringN(*args.args, C.int(args.arg_count))

	C.strcpy(error, C.CString("temporal() requires one argument"))

	// Rewrite the query for temporal tables
	rewrittenQuery, err := rewriteTemporalQuery(query)
	if err != nil {
		C.strcpy(error, C.CString(err.Error()))
		return nil
	}

	// Copy the rewritten query to the result
	result = C.CString(rewrittenQuery)
	*length = uint64(len(rewrittenQuery))

	return result
}

// Rewrite the query for temporal tables
func rewriteTemporalQuery(query string) (string, error) {
	query = strings.TrimSpace(query)
	lowerQuery := strings.ToLower(query)

	if strings.HasPrefix(lowerQuery, "insert") {
		return rewriteInsertQuery(query)
	} else if strings.HasPrefix(lowerQuery, "update") {
		return rewriteUpdateQuery(query)
	} else if strings.HasPrefix(lowerQuery, "delete") {
		return rewriteDeleteQuery(query)
	} else if strings.HasPrefix(lowerQuery, "select") {
		return rewriteSelectQuery(query)
	}

	return "", fmt.Errorf("unsupported query type")
}

// Rewrite an INSERT query to include historical data
func rewriteInsertQuery(query string) (string, error) {
	return query, nil
}

// Rewrite an UPDATE query to include historical data
func rewriteUpdateQuery(query string) (string, error) {
	// Extract the table name and condition from the update query
	tableName, condition := extractTableNameAndCondition(query)
	if tableName == "" {
		return "", fmt.Errorf("failed to extract table name")
	}

	// Rewrite the query to include the historical data handling
	rewrittenQuery := fmt.Sprintf(`
		INSERT INTO %s_history SELECT *, NOW() AS version_start, '9999-12-31 23:59:59' AS version_end FROM %s WHERE %s;
		UPDATE %s SET version_end = NOW() WHERE %s AND version_end = '9999-12-31 23:59:59';
		%s`,
		tableName, tableName, condition,
		tableName, condition,
		query)

	return rewrittenQuery, nil
}

// Rewrite a DELETE query to include historical data
func rewriteDeleteQuery(query string) (string, error) {
	// Extract the table name and condition from the delete query
	tableName, condition := extractTableNameAndCondition(query)
	if tableName == "" {
		return "", fmt.Errorf("failed to extract table name")
	}

	// Rewrite the query to include the historical data handling
	rewrittenQuery := fmt.Sprintf(`
		INSERT INTO %s_history SELECT *, NOW() AS version_start, '9999-12-31 23:59:59' AS version_end FROM %s WHERE %s;
		UPDATE %s SET version_end = NOW() WHERE %s AND version_end = '9999-12-31 23:59:59';
		%s`,
		tableName, tableName, condition,
		tableName, condition,
		query)

	return rewrittenQuery, nil
}

// Rewrite a SELECT query to include historical data
func rewriteSelectQuery(query string) (string, error) {
	// Check if the query requests historical data
	if strings.Contains(strings.ToLower(query), "for system_time") {
		tableName := extractTableName(query)
		if tableName == "" {
			return "", fmt.Errorf("failed to extract table name")
		}

		rewrittenQuery := fmt.Sprintf("SELECT * FROM %s_history WHERE %s", tableName, query)
		return rewrittenQuery, nil
	}

	return query, nil
}

// Extract table name and condition from a query
func extractTableNameAndCondition(query string) (string, string) {
	lowerQuery := strings.ToLower(query)
	fromIndex := strings.Index(lowerQuery, " from ")
	if fromIndex == -1 {
		return "", ""
	}
	fromIndex += len(" from ")

	// Find end of table name
	endIndex := fromIndex
	for endIndex < len(query) && query[endIndex] != ' ' {
		endIndex++
	}
	tableName := query[fromIndex:endIndex]

	// Extract the condition
	whereIndex := strings.Index(lowerQuery, " where ")
	if whereIndex == -1 {
		return tableName, "1=1" // No condition, match all rows
	}
	whereIndex += len(" where ")
	condition := query[whereIndex:]

	return tableName, condition
}

// Extract table name from a query
func extractTableName(query string) string {
	lowerQuery := strings.ToLower(query)
	fromIndex := strings.Index(lowerQuery, " from ")
	if fromIndex == -1 {
		return ""
	}
	fromIndex += len(" from ")

	// Find end of table name
	endIndex := fromIndex
	for endIndex < len(query) && query[endIndex] != ' ' {
		endIndex++
	}
	return query[fromIndex:endIndex]
}

func main() {}

// The C main function to build the shared library
/*
int main() {
    return 0;
}
*/
