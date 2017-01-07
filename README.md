## Lytics Command Line tool & Developers Aid

The goal of this tool is to provide CLI access to Lytics api, 
as well as a developers aid to enable writing and testing LQL 
as easily as possible.

Would love any feature requests or ideas that would make this useful to you.

### Installation

Download a binary from https://github.com/lytics/lytics/releases and rename to lytics.
Or install from source.

```
# Or if you have go installed
go get -u github.com/lytics/lytics

```

### Usage

All examples use JQ https://stedolan.github.io/jq/ to 
prettify the json output.

```

export LIOKEY="your_api_key"
lytics --help


```

**Segment Scan Usage**

Exporting CSV files, with usage.

```

# Scan a segment by id
lytics --id=ab93a9801a72871d689342556b0de2e9 segmentscan | jq '.'

# Scan a segment by Slug
lytics --id=last_2_hours segmentscan | jq '.'

# write out this segment to temp file so we can play with jq

lytics --id=last_2_hours segmentscan > /tmp/users.json

# same thing but with "Ad hoc query"
lytics segmentscan '

FILTER AND (
    lastvisit_ts > "now-2d"
    EXISTS email
)
FROM user

' > /tmp/users.json


# use JQ to output a few fields

cat /tmp/users.json | \
 jq -c ' {country: .country, city: .city, org: .org, uid: ._uid, visitct: .visitct} '

# create a csv file from these users
echo "country,city,org,uid,visitct\n" > /tmp/users.csv
cat /tmp/users.json | \
 jq -r ' [ .country, .city, .org,  ._uid, .visitct ] | @csv ' >> /tmp/users.csv

```


**Lytics Watch Usage**

1.  Create NAME.lql (any name) file in a folder.
2.  Assuming you already have data collected, it will use our api to show recent examples against that lql.

You can open and edit in an editor.  Every time you edit it will print resulting
users it interpreted from recent data to our api.


### Example

```

# get your api key from web app account settings
export LIOKEY="your_api_key"

cd /path/to/your/project

# create an lql file
# - utilize the lytics app "data -> streams" section to see
#   data fields you are sending to lytics.

# you can create this in an editor as well
echo '
SELECT 
   user_id,
   name,
   todate(ts),
   match("user.") AS user_attributes,
   map(event, todate(ts))   as event_times   KIND map[string]time  MERGEOP LATEST

FROM default
INTO USER
BY user_id
ALIAS my_query
' > default.lql


# start watching
lytics watch | jq '.'

# now edit, json results of how data is interpreted is output


```



**Lytics Watch With Custom Data**
1.  Create NAME.lql (any name) file in a folder.
2.  Create NAME.json (any name, must match lql file name) in folder.
3.  run the `lytics watch` command from the folder with files.
4.  Edit .lql, or .json files, upon change the evaluated result of .lql, json will be output.


### Example

```

# get your api key from web app account settings
export LIOKEY="your_api_key"

cd /tmp

# start watching in background
lytics watch &


# create an lql file
echo '
SELECT 
   user_id,
   name,
   todate(ts),
   match("user.") AS user_attributes,
   map(event, todate(ts))   as event_times   KIND map[string]time  MERGEOP LATEST

FROM data
INTO USER
BY user_id
ALIAS hello
' > hello.lql

# Create a Json data of events to feed into lql query
echo '[
    {"user_id":"dump123","name":"Down With","company":"Trump", "event":"project.create", "ts":"2016-11-09"},
    {"user_id":"another234","name":"No More","company":"Trump", "event":"project.signup","user.city":"Portland","user.state":"Or", "ts":"2016-11-09"}
]' > hello.json



```
