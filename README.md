# execql
Execute multiple CQL from file with concurrency.  

# Install
with golang  
```
$ go get -u github.com/ken-aio/execql
```

Mac or Linux  
Below sample is linux. please see releases.  
```
$ wget https://github.com/ken-aio/execql/releases/download/v0.0.1/execql_v0.0.1_linux_amd64.tar.gz
$ tar xvf execql_v0.0.1_linux_amd64.tar.gz
$ mv execql /usr/local/bin/
```

# Usage
```
Usage:
  execql [flags]

Flags:
  -f, --file string       cql file path (required)
  -h, --help              help for execql
  -H, --host string       cassandra host. split ',' if many host. e.g.) cassandra01, cassandra02
 (default "localhost")
  -k, --keyspace string   exec target keyspace (required)
  -n, --num-conns int     connection nums (default 10)
  -p, --password string   connection password
  -P, --port int          cassandra port (default 9042)
  -t, --thread int        concurrent query request thread num (default 1)
      --timeout int       query timeout(ms) (default 60000)
  -u, --user string       connection user
```

## sample command
```
$ execql -k test-keyspace -f /path/to/exec.cql -n 10 -t 10
2019/01/21 16:15:06 Reading input cql file... /path/to/exec.cql
2019/01/21 16:15:06 Complete reading input cql file
2019/01/21 16:15:06 Creating cassandra session...
2019/01/21 16:15:16 Complete creating cassandra session
2019/01/21 16:15:16 Execute CQL...
2019/01/21 22:00:00 start thread#9      / cql num is 10573
2019/01/21 22:00:00 start thread#3      / cql num is 10578
2019/01/21 22:00:00 start thread#0      / cql num is 10578
2019/01/21 22:00:00 start thread#1      / cql num is 10578
2019/01/21 22:00:00 start thread#6      / cql num is 10578
2019/01/21 22:00:00 start thread#2      / cql num is 10578
2019/01/21 22:00:00 start thread#5      / cql num is 10578
2019/01/21 22:00:00 start thread#4      / cql num is 10578
2019/01/21 22:00:00 start thread#7      / cql num is 10578
2019/01/21 22:00:00 start thread#8      / cql num is 10578
2019/01/21 22:01:46 Complete thread#7
2019/01/21 22:01:47 Complete thread#9
2019/01/21 22:01:47 Complete thread#0
2019/01/21 22:01:47 Complete thread#6
2019/01/21 22:01:47 Complete thread#4
2019/01/21 22:01:47 Complete thread#3
2019/01/21 22:01:48 Complete thread#8
2019/01/21 22:01:48 Complete thread#5
2019/01/21 22:01:48 Complete thread#1
2019/01/21 22:01:48 Complete thread#2
2019/01/21 16:15:16 Complete execute CQL
```
