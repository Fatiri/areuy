# AREUY Times

# Instatlations
```sh
$ go get github.com/Fatiri/areuy/times
```

## List Of interfaces of Time

```sh
type Time interface {
	Now(timeGMT *int) time.Time
	TimeStampToDateStr(timeStr, layout string) string
	TimeStampToDate(timeStr, layout string) time.Time
}
```

- Now() 

```
    return currentime in GMT Zone
```

- TimeStampToDateStr() 

 ```
    will manipulate timestamp and return the date string

    Parameter : 
    1. timeStr 
       time with string data type will manipulated
    2. layout
       layout is the date layout using for 
 ```

 
- TimeStampToDate(timeStr, layout string) time.Time


```
    will manipulate timestamp and return the date

    Parameter : 
    1. timeStr 
       time with string data type will manipulated
    2. layout
       layout is the date layout using for 
 ```
