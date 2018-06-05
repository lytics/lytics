package cmds

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/apcera/termtables"
	"github.com/araddon/gou"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

var (
	apikey       string
	outputFormat string

/*
	// TODO DEPRECATE
	dataapikey string
	method     string
	id         string
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
*/
)

var (
	app      *cli.App
	client   *lytics.Client
	commands = make([]cli.Command, 0)
)

func init() {
	gou.SetupLogging("debug")
	gou.SetColorOutput()
}
func addCommand(c cli.Command) {
	commands = append(commands, c)
}
func Run() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app = cli.NewApp()
	app.Usage = "Lytics command line tools"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "authtoken, t",
			Usage:       "auth token for Lytics api",
			EnvVar:      "LIOKEY,LYTICS_AUTH_TOKEN",
			Destination: &apikey,
		},
		cli.StringFlag{
			Name:        "format, f",
			Usage:       "Format [json, table] to print results as",
			Value:       "table",
			Destination: &outputFormat,
		},
	}
	app.Before = func(c *cli.Context) error {
		client = lytics.NewLytics(apikey, nil)
		return nil
	}
	app.Commands = commands

	err := app.Run(os.Args)
	exitIfErr(err, "Could not run")
}

func resultWrite(cliCtx *cli.Context, result interface{}) {
	switch outputFormat {
	case "table":
		switch val := result.(type) {
		case []lytics.TableWriter:
			resultWriteTable(cliCtx, val)
		case lytics.TableWriter:
			resultWriteTable(cliCtx, []lytics.TableWriter{val})
		default:
			exitIfErr(fmt.Errorf("expected tablewriter got %T", result), "Wrong type")
		}

	case "json":
		jsonOut, err := json.MarshalIndent(result, "", "  ")
		exitIfErr(err, "Could not marshal json")
		fmt.Printf("%s\n", string(jsonOut))
	}
}

func resultWriteTable(cliCtx *cli.Context, list []lytics.TableWriter) {
	table := termtables.CreateTable()
	for i, row := range list {
		if i == 0 {
			table.AddHeaders(row.Headers()...)
		}
		table.AddRow(row.Row()...)
	}
	fmt.Println(table.Render())
}

/*
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
*/

func exitIfErr(err error, msg string, args ...interface{}) {
	if err != nil {
		args = append(args, err)
		fmt.Fprintf(os.Stderr, msg+"err=%v\n", args)
		os.Exit(1)
	}
}

func errExit(err error, msg string) {
	fmt.Fprintf(os.Stderr, "%v: %s\n", err, msg)
	os.Exit(1)
}
