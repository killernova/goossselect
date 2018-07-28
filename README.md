# goossselect
Ali OSS select package for go

## INSTALL
`go get -u github.com/killernova/goossselect`

## USAGE

### Step 0: import package
`import(gol github.com/killernova/goossselect)`

### Step 1: create a new selector or meta
`s, err := gol.NewSelector("Select * from ossobject")`

You can also set options for this generator with predefined functions or functions defined by yourself:

`s, err := gol.NewSelector("Select * from ossobject", gol.SetRange("line", 10, 333))`

### Step 2: create a new config object
```go
c := gol.NewConfig(func(conf *gol.Config) {
		conf.AccessKeyID = "ercbcv43tgert"
		conf.AccessKeySecret = "sferg456ghvbsf435435"
	})
```

### Step 3: call function
```go
data, err = s.SelectQuery("dwl-oss", "links.csv", c)
```
or you can save data to file
```go
f, err := os.OpenFile("efg.csv", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = s.SelectQueryToFile("dwl-oss", "links.csv", c, f)
	if err != nil {
		fmt.Println(err)
	}
```
maybe you want query meta data
```go
m, err := gol.NewMeta()
data, err = m.SelectMeta("dwl-oss", "links.csv", c)
```
this will return a struct which contains informations of your query csv file
```go
type MetaResponse struct {
	Lines int64
	Columns int64
	Splits int64
	ContentLength int64
}
```

## OPTIONS

Most of time you can use the default options, but if you want to change it, use predefined functions as follows:
```go
# Selector:
SetFileHeaderInfo(s string)
SetInRecordDelimiter(s string)
SetInFieldDelimiter(s string)
SetOutRecordDelimiter(s string)
SetOutFieldDelimiter(s string)
SetOutputRawData(b bool)
SetKeepAllColumns(b bool)
SetQuoteCharacter(s string)
SetCommentCharacter(s string)
SetRange(s string, start, end int)

# Meta:
SetMetaOverwriteIfExisting(b bool)
SetMetaRecordDelimiter(s string)
SetMetaFieldDelimiter(s string)
SetMetaQuoteCharacter(s string)
```

You can also use your own functions which receive *Selector or *Meta and return error

