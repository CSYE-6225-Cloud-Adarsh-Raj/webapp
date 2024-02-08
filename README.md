# webapp
A Golang web application to authenticate and store users with PostgresQL

I.  Prerequisites for building and deploying your application locally:
    Download go from its official website (go version go1.21.6)
    
    Extract go:
    > sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

    Set Up Path:
    > export PATH=$PATH:/usr/local/go/bin

II. Build and Deploy instructions for the web application:
    
    Clone the org repo and cd into the project
    Run following command:
        > go mod tidy //installs dependencies
    Run following command:
        > go build .
        > go run main.go

    Database Operrations:
    Commands for Postgres:
    > PG_CTL is used to start,stop server
    > psql is the postgres cli used to query the database

    To login PSQL:
    > sudo -u postgres psql(username - postgres)

    Start server:
    > sudo -u postgres pg_ctl -D /Library/PostgreSQL/16/data start
    > pgctl start (alias set in .zshrc)

    Stop Server:
    > pgctl stop

    Status of Server:
    > pgctl status

    Get inside database through psql:
    > psql -U postgres -d <database-name>
    > psql -U postgres -d test

    List contents inside table:
    > SELECT * FROM <table-name>
    > SELECT * FROM user_models;

III. REST APIs endpoint:
    
    Availaible in swagger docs. Link - https://app.swaggerhub.com/apis-docs/csye6225-webapp/cloud-native-webapp/2024.spring.02


GIT Workflow:

    Clone the  Developers Fork
    On your Developers Fork branch:
        > After comming changes to <dev_branch_name_fork>
        > git push fork <dev_branch_name_fork>
        > Create a PR from your <dev_branch_name_fork> to main branch of ORG repo
    Merge PR
    Delete the  <dev_branch_name_fork>
    
    Sync Fork:
        >On local git, checkout to fork main
        >git pull upstream main
        >git push fork main

References:

Assignemnt 1:
    
    https://medium.com/@venu-prasanna/developing-a-restful-api-with-go-gin-and-gorm-part-1-router-setup-db-configuration-a31a74ad416d

Assignment 2:

    https://www.alexedwards.net/blog/basic-authentication-in-go

    https://medium.com/@venu-prasanna/developing-a-restful-api-with-go-gin-and-gorm-part-2-repository-setup-table-driven-testing-7d18cc532b65

    https://gorm.io/docs/index.html

    https://www.postgresql.org/docs/


On CentOS 8
Download postgres 16
> sudo dnf -y install https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm

Disable old version of postgres
> sudo dnf -qy module disable postgresql

Install
> sudo dnf -y install postgresql16-server

Initialise the Database
> sudo /usr/pgsql-16/bin/postgresql-16-setup initdb

Enable the service
> sudo systemctl enable postgresql-16

Start the service
> sudo systemctl start postgresql-16

Verify Installation:
> sudo -u postgres psql -c "SELECT version();"

Setup Golang
> curl -LO https://golang.org/dl/go1.21.6.linux-386.tar.gz
> sudo tar -C /usr/local -xzf go1.21.6.linux-386.tar.gz
> echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bash_profile
> source ~/.bash_profile
> go version

Set ENV
> echo 'export DB_HOST=localhost' >> ~/.bash_profile
> echo 'export DB_USER=test' >> ~/.bash_profile
> echo 'export DB_PASSWORD=test' >> ~/.bash_profile
> echo 'export DB_NAME=test' >> ~/.bash_profile
> source ~/.bash_profile

To unzip the project
> Install unzip tool
> sudo dnf install unzip

Unzip the project
> unzip project.zip

PSQl commands:
> sudo -u postgres psql

Inside psql cli
> CREATE USER test WITH ENCRYPTED PASSWORD 'test';
> CREATE DATABASE test;

> GRANT ALL PRIVILEGES ON DATABASE test TO test;
-- Grant usage on the schema
> GRANT USAGE ON SCHEMA public TO test;

-- Grant create on the schema
> GRANT CREATE ON SCHEMA public TO test;

-- Grant all privileges on all tables in the public schema
> GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO test;

-- Optionally, grant privileges on sequences
> GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO test;
