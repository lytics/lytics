module lytics

go 1.16

require (
	github.com/araddon/gou v0.0.0-20190110011759-c797efecbb61
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/lytics/go-lytics v0.0.0-20210216161913-1f2ffd155680
	github.com/lytics/lytics v0.0.0-20210130003108-a19f4c39c6c9
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
)

replace github.com/lytics/lytics => ../lytics
