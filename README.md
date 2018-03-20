# go-microservice
playground for GO programming


------------------------------------
Database:
```shell
sudo docker run --name posttest -p 5432:5432 -e POSTGRES_PASSWORD=fred postgres:alpine
```

table:
```sql
  create table users(
    id  serial PRIMARY KEY,
    name varchar(256) NOT NULL,
    email varchar(100)  NOT NULL,
    password varchar(256)  NOT NULL
  )
```

Using different ways how to handle database
* viaGlobalVariable - db connection is kept in global variable and reused by models
* viaIncjection - db connection is past as a parameter to model
* viaInterface - db connection is hidden behind interface and models are past as a parameter

based on:
http://www.alexedwards.net/blog/organising-database-access
