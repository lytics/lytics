package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	lytics "github.com/lytics/go-lytics"
)

var (
	apikey        string
	dataapikey    string
	method        string
	id            string
	output        string
	file          string
	fields        string
	fieldsSlice   []string
	segments      string
	segmentsSlice []string
	entitytype    string
	fieldname     string
	fieldvalue    string
	table         string
	limit         int
)

type Cli struct {
	Client *lytics.Client
}

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		usageExit()
	}

	flag.StringVar(&apikey, "apikey", os.Getenv("LIOKEY"), "Lytics API Key - Or use env LIOKEY")
	flag.StringVar(&dataapikey, "dataapikey", os.Getenv("LIODATAKEY"), "Lytics Data API Key - Or use env LIODATAKEY")
	flag.StringVar(&id, "id", "", "Id of object")
	flag.StringVar(&segments, "segments", "", "Comma Separated Segments")
	flag.StringVar(&fields, "fields", "", "Comma Separated Fields")
	flag.StringVar(&fieldname, "fieldname", "", "Field Name")
	flag.StringVar(&fieldvalue, "fieldvalue", "", "Field Value")
	flag.StringVar(&entitytype, "entitytype", "", "Entity Type")
	flag.StringVar(&table, "table", "", "Schema Table")
	flag.StringVar(&file, "file", "", "Output File Name")
	flag.IntVar(&limit, "limit", -1, "Result Limit")
	flag.Parse()
}

func main() {
	if apikey == "" && dataapikey == "" {
		fmt.Println(`Missing -apikey and/or -method: use -help for assistance

	LIOKEY env variable will fullfill api key needs
	`)
		os.Exit(1)
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}
	method = flag.Args()[0]

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create lytics client with auth info
	c := Cli{
		Client: lytics.NewLytics(apikey, dataapikey, nil),
	}

	output, err := c.handleFunction(method)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}

	if file != "" {
		err := writeToFile(file, output)
		if err != nil {
			fmt.Println("Failed to write data to file %s: %v", file, err)
			return
		}

		fmt.Println("Data written to %s successfully.", file)
		return
	}

	if len(output) > 0 {
		fmt.Println(output)
	}
}

func writeToFile(file, data string) error {
	err := ioutil.WriteFile(file, []byte(data), 0644)
	return err
}

func appendToFile(file, data string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		return err
	}

	return nil
}

func (c *Cli) handleFunction(method string) (string, error) {
	var (
		result interface{}
		err    error
	)

	if fields != "" {
		fieldsSlice = strings.Split(fields, ",")
	}

	if segments != "" {
		segmentsSlice = strings.Split(segments, ",")
	}

	switch method {
	case "account":
		result, err = c.getAccounts(id)

	case "auth":
		result, err = c.getAuths(id)

	case "schema":
		result, err = c.getSchema(table)

	case "entity":
		result, err = c.getEntity(entitytype, fieldname, fieldvalue, fieldsSlice)

	case "provider":
		result, err = c.getProviders(id)

	case "segment":
		result, err = c.getSegments("user", segmentsSlice)

	case "segmentscan":
		if id == "" && len(flag.Args()) == 2 {
			id = flag.Args()[1]
		}
		c.getEntityScan(id, limit, func(e *lytics.Entity) {
			fmt.Println(e.PrettyJson())
		})
		return "", nil

	case "segmentsize":
		result, err = c.getSegmentSizes(segmentsSlice)

	case "segmentattribution":
		result, err = c.getSegmentAttributions(segmentsSlice, limit)

	case "work":
		result, err = c.getWorks(id)

	case "user":
		result, err = c.getUsers(id)

	case "query":
		result, err = c.getQueries(id)

	case "watch":
		c.watch()
	default:
		flag.Usage()
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return makeJSON(result), err
}

func validate() bool {
	return true
}

func makeJSON(blob interface{}) string {
	jsonOut, err := json.MarshalIndent(blob, "", "	")
	if err != nil {
		return fmt.Sprintf("Failed: %v", err)
	}

	return string(jsonOut)
}

func exitIfErr(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %s\n", err, msg)
		os.Exit(1)
	}
}

func errExit(err error, msg string) {
	fmt.Fprintf(os.Stderr, "%v: %s\n", err, msg)
	os.Exit(1)
}

func usageExit() {
	fmt.Printf(`
--------------------------------------------------------
**************  LYTICS COMMAND LINE HELP  **************
--------------------------------------------------------

ENV Vars:
    export LIOKEY="your_api_key"
    export LIODATAKEY="your_api_key"

METHODS:
    [account]
         retrieves account information based upon api key.
         if no id is passed, all accounts returned.
         -------
         params:
         -------
              <id>               OPTIONAL       string
         -------
         example:
         -------
              lytics account
              lytics --id=<id> account

    [schema]
         retrieves table schema (fields, types)
         -------
         params:
         -------
              <table>            REQUIRED       string
              <limit>            OPTIONAL       int
         -------
         example:
         -------
              lytics schema
              lytics --limit=<limit> --table=user schema

    [entity]
         retrieves entity (a single user) information
         -------
         params:
         -------
              <entitytype>       REQUIRED       string (user or content)
              <fieldname>        REQUIRED       string (name of field to search by, e.g. email)
              <fieldvalue>       REQUIRED       string (value of field to search by)
              <fields>           OPTIONAL       string (comma separated list of fields to filter by)
         -------
         example:
         -------
              lytics -entitytype=user -fieldname=email -fieldvalue="me@me.com" entity
              lytics -entitytype=user -fieldname=email -fieldvalue="me@me.com" -fields=email entity

    [segmentscan]
         retrieves a list of users (actually entities, could be content, etc).
         
         -------
         params:
         -------
              <id>   id=id_or_slug
         -------
         example:
         -------
              lytics --id=slug_of_segment segmentscan

              # use a segment QL query
              lytics segmentscan '
                  FILTER AND (
                     EXISTS email 
                     last_active_ts > "now-7d"
                  )
              '
              
              # see what single user looks like
              lytics --limit=1 segmentscan ' FILTER * FROM user'

              # see what content looks like 
              lytics --limit=1 segmentscan ' FILTER * FROM content'

    [segment]
         retrieves segment information based upon api key.
         if no id is passed, all segments returned.
         -------
         params:
         -------
              <segments>         OPTIONAL       string (comma separated list of segment ids, max 1)
         -------
         example:
         -------
              lytics segment
              lytics -segments=slug_of_segment segment

    [segmentsize]
         retrieves segment sizes information based upon api key.
         if no id is passed, all segment sizes returned.
         -------
         params:
         -------
              <segments>         OPTIONAL       string (comma separated list of segment ids)
         -------
         example:
         -------
              lytics segmentsize
              lytics -segmentes=one,two segmentsize

    [segmentattribution]
         retrieves segment information based upon api key.
         if no id is passed, all segments returned.
         -------
         params:
         -------
              <segments>         OPTIONAL       string (comma separated list of segment ids)
              <limit>            OPTIONAL       int
         -------
         example:
         -------
              lytics segmentattribution
              lytics -segments=one,two -limit=5 segmentattribution

    [user]
         retrieves administrative user information based upon api key.
         if no id is passed, all users returned.
         -------
         params:
         -------
              <id>               OPTIONAL       string
         -------
         example:
         -------
              lytics user
              lytics -id=<id> user

    [query]
         retrieves query information
         -------
         params:
         -------
              <alias>               OPTIONAL       string
         -------
         example:
         -------
              lytics query
              lytics --id=<alias> query

    [watch]
         watch the current folder for .lql, .json files to evaluate
         the .lql query against the data in .json to preview output.

         .lql file name must match the json file name.

         For Example: 
            cd /tmp 
            ls *.lql       # assume a temp.lql
            cat temp.json  # file of data

         -------
         example:
         -------
              lytics watch

`)
	os.Exit(1)
}
