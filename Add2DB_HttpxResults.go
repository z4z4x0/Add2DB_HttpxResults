package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	Timestamp        string                 `json:"timestamp"`
	Asn              map[string]interface{} `json:"asn,omitempty"`
	Tls              map[string]interface{} `json:"tls,omitempty"`
	Hash             map[string]interface{} `json:"hash,omitempty"`
	CdnName          string                 `json:"cdn_name,omitempty"`
	Port             string                 `json:"port,omitempty"`
	URL              string                 `json:"url"`
	Input            string                 `json:"input"`
	Title            string                 `json:"title,omitempty"`
	Scheme           string                 `json:"scheme,omitempty"`
	ContentType      string                 `json:"content_type,omitempty"`
	Method           string                 `json:"method,omitempty"`
	Host             string                 `json:"host,omitempty"`
	Path             string                 `json:"path,omitempty"`
	Favicon          string                 `json:"favicon,omitempty"`
	FaviconPath      string                 `json:"favicon_path,omitempty"`
	FinalURL         string                 `json:"final_url,omitempty"`
	Time             string                 `json:"time,omitempty"`
	ChainStatusCodes []int                  `json:"chain_status_codes,omitempty"`
	A                []string               `json:"a,omitempty"`
	Cname            []string               `json:"cname,omitempty"`
	Tech             []string               `json:"tech,omitempty"`
	Words            int                    `json:"words,omitempty"`
	Lines            int                    `json:"lines,omitempty"`
	StatusCode       int                    `json:"status_code,omitempty"`
	ContentLength    int                    `json:"content_length,omitempty"`
	Failed           bool                   `json:"failed"`
	Vhost            bool                   `json:"vhost,omitempty"`
	Cdn              bool                   `json:"cdn,omitempty"`
	Knowledgebase    map[string]interface{} `json:"knowledgebase,omitempty"`
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS records (
        timestamp TEXT,
        asn TEXT,
        tls TEXT,
        hash TEXT,
        cdn_name TEXT,
        port TEXT,
        url TEXT,
        input TEXT,
        title TEXT,
        scheme TEXT,
        content_type TEXT,
        method TEXT,
        host TEXT,
        path TEXT,
        favicon TEXT,
        favicon_path TEXT,
        final_url TEXT,
        time TEXT,
        chain_status_codes TEXT,
        a TEXT,
        cname TEXT,
        tech TEXT,
        words INTEGER,
        lines INTEGER,
        status_code INTEGER,
        content_length INTEGER,
        failed BOOLEAN,
        vhost BOOLEAN,
        cdn BOOLEAN,
        knowledgebase TEXT,
        UNIQUE(timestamp, asn, tls, hash, cdn_name, port, url, input, title, scheme, content_type, method, host, path, favicon, favicon_path, final_url, time, chain_status_codes, a, cname, tech, words, lines, status_code, content_length, failed, vhost, cdn, knowledgebase) ON CONFLICT IGNORE
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("./data")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			processFile(filepath.Join("./data", file.Name()), db)
		}
	}
}

func processFile(fileName string, db *sql.DB) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rec Record
		err := json.Unmarshal(scanner.Bytes(), &rec)
		if err != nil {
			log.Println("Error parsing JSON: ", err)
			continue
		}

		values, placeholders := prepareInsertValues(rec)
		insertSQL := fmt.Sprintf(`INSERT INTO records (%s) VALUES (%s)`, placeholders, makeQuestionMarks(values))
		_, err = db.Exec(insertSQL, values...)
		if err != nil {
			log.Println("Error inserting record: ", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file:", err)
	}
}

// prepareInsertValues prepares the field names and values for insertion.
func prepareInsertValues(rec Record) ([]interface{}, string) {
	v := reflect.ValueOf(rec)
	typeOfS := v.Type()
	var fieldNames []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		jsonTag := field.Tag.Get("json")
		fieldName := strings.Split(jsonTag, ",")[0] // Correctly handle JSON tags
		fieldNames = append(fieldNames, fieldName)

		value := v.Field(i).Interface()
		if mapVal, ok := value.(map[string]interface{}); ok {
			if jsonStr, err := json.Marshal(mapVal); err == nil {
				value = string(jsonStr)
			}
		} else if sliceVal, ok := value.([]int); ok {
			if jsonStr, err := json.Marshal(sliceVal); err == nil {
				value = string(jsonStr)
			}
		} else if sliceVal, ok := value.([]string); ok {
			if jsonStr, err := json.Marshal(sliceVal); err == nil {
				value = string(jsonStr)
			}
		}
		values = append(values, value)
	}

	return values, fmt.Sprintf(`"%s"`, strings.Join(fieldNames, `", "`))
}

// makeQuestionMarks generates the placeholder string for SQL insertion.
func makeQuestionMarks(values []interface{}) string {
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}
