package cmds

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/araddon/gou"
	lytics "github.com/lytics/go-lytics"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var (
	apikey          string
	outputFormat    string
	userCreatedName string
)

var (
	app      *cli.App
	client   *lytics.Client
	commands = make([]*cli.Command, 0)
)

func init() {
	gou.SetupLogging("debug")
	gou.SetColorOutput()
}
func addCommand(c cli.Command) {
	commands = append(commands, &c)
}

// Run main entrypoint for CLI command.
func Run() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app = cli.NewApp()
	app.Usage = "Lytics command line tools"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "authtoken, t",
			Usage:       "Auth token for Lytics API",
			EnvVars:     []string{"LIOKEY", "LYTICS_AUTH_TOKEN"},
			Destination: &apikey,
		},
		&cli.StringFlag{
			Name:        "format, f",
			Usage:       "Format [json, table, csv] to print results as",
			Value:       "table",
			Destination: &outputFormat,
		},
		&cli.StringFlag{
			Name:        "name, n",
			Usage:       "Name for CSV filename",
			Value:       "",
			Destination: &userCreatedName,
		},
	}
	app.Before = func(c *cli.Context) error {
		client = lytics.NewLytics(apikey, nil)
		return nil
	}
	app.Commands = commands

	err := app.Run(os.Args)
	exitIfErr(err, "could not run")
}

func resultWrite(cliCtx *cli.Context, result interface{}, name string) {
	if userCreatedName != "" {
		name = userCreatedName
	}

	switch outputFormat {
	case "table":
		switch val := result.(type) {
		case []lytics.TableWriter:
			resultWriteTable(cliCtx, val)
		case lytics.TableWriter:
			resultWriteTable(cliCtx, []lytics.TableWriter{val})
		default:
			exitIfErr(fmt.Errorf("expected tablewriter got %T", result), "wrong type")
		}

	case "json":
		jsonOut, err := json.MarshalIndent(result, "", "  ")
		exitIfErr(err, "could not marshal JSON")
		fmt.Printf("%s\n", string(jsonOut))

	case "csv":
		switch val := result.(type) {
		case []lytics.TableWriter:
			resultWriteCSV(val, name)
		case lytics.TableWriter:
			resultWriteCSV([]lytics.TableWriter{val}, name)
		default:
			exitIfErr(fmt.Errorf("expected tablewriter got %T", result), "wrong type")
		}
	}

}

func resultWriteCSV(list []lytics.TableWriter, name string) {
	filename := fmt.Sprintf("%s.csv", name)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	h := toStringArr(list[0].Headers())
	if err := w.Write(h); err != nil {
		log.Fatalln("Error writing header to CSV:", err)
	}

	for _, item := range list {
		r := toStringArr(item.Row())
		if err := w.Write(r); err != nil {
			log.Fatalln("Error writing record to CSV:", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func resultWriteTable(cliCtx *cli.Context, list []lytics.TableWriter) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	for i, row := range list {
		if i == 0 {
			table.SetHeader(rowToString(row.Headers()))
		}
		table.Append(rowToString(row.Row()))
	}
	table.SetAutoFormatHeaders(false)
	table.Render()
	fmt.Println(tableString.String())
}

func rowToString(row []interface{}) []string {
	res := make([]string, 0, len(row))
	for _, val := range row {
		str, ok := val.(string)
		if ok {
			res = append(res, str)
			continue
		}
		res = append(res, fmt.Sprint(val))
	}

	return res
}

func exitIfErr(err error, msg string, args ...interface{}) {
	if err != nil {
		args = append(args, err)
		fmt.Fprintf(os.Stderr, msg+" err=%v\n", args)
		os.Exit(1)
	}
}

func errExit(err error, msg string) {
	fmt.Fprintf(os.Stderr, "%v: %s\n", err, msg)
	os.Exit(1)
}

func toStringArr(inter []interface{}) []string {
	strArr := make([]string, len(inter))
	for i, field := range inter {
		strArr[i] = fmt.Sprint(field)
	}
	return strArr
}
