# genlog

genlog is MySQL General Query Log parser.

## Usage

```
$ cat general.log
2020-05-27T05:03:27.500301Z   11 Query	SET @@sql_log_bin=off
2020-05-27T05:03:27.543379Z   11 Query	select @@session.tx_read_only
2020-05-27T05:03:27.683485Z   11 Query	COMMIT
...

$ genlog general.log # or `cat general.log | genlog`
{"Time":"2020-05-27T05:03:27.500301Z","Id":"11","Command":"Query","Argument":"SET @@sql_log_bin=off"}
{"Time":"2020-05-27T05:03:27.543379Z","Id":"11","Command":"Query","Argument":"select @@session.tx_read_only"}
{"Time":"2020-05-27T05:03:27.683485Z","Id":"11","Command":"Query","Argument":"COMMIT"}
...
```
