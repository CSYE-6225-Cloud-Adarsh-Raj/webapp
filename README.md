# webapp
A Golang web application to authenticate and store users with PostgresQL

Commands for Postgres:
PG_CTL is used to start,stop server
psql is the postgres cli used to query the database

To login PSQL:
> sudo -u postgres psql(username - postgres)

> Start server:
1. sudo -u postgres pg_ctl -D /Library/PostgreSQL/16/data start
2. pgctl start (alias set in .zshrc)

> Stop Server:
pgctl stop

> Status of Server:
pgctl status

Get inside database through psql:
psql -U postgres -d <database-name>
psql -U postgres -d test

List contents inside table:
SELECT * FROM <table-name>
SELECT * FROM user_models;


I>  Prerequisites for building and deploying your application locally:
    i. Clone the repo and cd into the project
    ii. Run following command:
        go mod tidy //installs dependencies

II> Build and Deploy instructions for the web application:
    i. Run following command:
        go build .
        go run main.go

III> REST APIs endpoint:
    Availaible in swagger docs. Link - https://app.swaggerhub.com/apis-docs/csye6225-webapp/cloud-native-webapp/2024.spring.02


GIT Workflow:
Clone the  Developers Fork
On your Developers Fork branch:
    ...After comming changes to <dev_branch_name_fork>
    git push fork <dev_branch_name_fork>
    Create a PR from your <dev_branch_name_fork> to main branch of ORG repo
Merge PR
Delete the  <dev_branch_name_fork>

Sync Fork:
    On local git, checkout to fork main
    git pull upstream main
    git push fork main

References:
Assignemnt 1:

https://medium.com/@venu-prasanna/developing-a-restful-api-with-go-gin-and-gorm-part-1-router-setup-db-configuration-a31a74ad416d

Assignment 2:

1. https://www.alexedwards.net/blog/basic-authentication-in-go
2. https://medium.com/@venu-prasanna/developing-a-restful-api-with-go-gin-and-gorm-part-2-repository-setup-table-driven-testing-7d18cc532b65
