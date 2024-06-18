## About

MySQL user defined function in Go.

## Dependencies

- Mysql: sudo apt install mysql-server
- libmysqlclient: sudo apt install libmysqlclient-dev -y

## Run MySQL

```console
docker run --name mysql -e MYSQL_ROOT_PASSWORD=password -p 3306:3306 -d mysql
```

## Install

```console
make install
```

## Use the function

```sql
SELECT temporal(1);

mysql> select temporal(1);
+-------------+
| temporal(1) |
+-------------+
| hello world |
+-------------+
1 row in set (0.00 sec)
```

## Finding the mysql include dir

```console
mysql_config --include
# -I/usr/include/mysql
```

## Finding the mysql plugin dir

```console
SELECT @@plugin_dir;
```