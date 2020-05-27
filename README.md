# genlog

genlog is MySQL General Query Log parser.

## Usage

```sh
$ genlog general.log
{"Time":"2020-05-27T05:03:27.500301Z","Id":"11","Command":"Query","Argument":"SET @@sql_log_bin=off\n"}
{"Time":"2020-05-27T05:03:27.543379Z","Id":"11","Command":"Query","Argument":"select @@session.tx_read_only\n"}
{"Time":"2020-05-27T05:03:27.683485Z","Id":"11","Command":"Query","Argument":"COMMIT\n"}
...
```
