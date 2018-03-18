# go-microservice
playground for GO programming


------------------------------------
Database:
sudo docker run --name posttest -p 5432:5432 -e POSTGRES_PASSWORD=fred postgres:alpine


table:
  create table users(
    id  serial PRIMARY KEY,
    name varchar(256) NOT NULL,
    email varchar(100)  NOT NULL,
    password varchar(256)  NOT NULL
  )
------------------------------------
