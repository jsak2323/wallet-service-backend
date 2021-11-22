## Install VM Requirements
    export LANGUAGE=en_US.UTF-8
    export LC_ALL=en_US.UTF-8

    sudo apt-get update
    sudo apt-get -y upgrade
    sudo apt-get install mysql-server 
    sudo apt-get install software-properties-common
    sudo apt-get update && sudo apt-get install sqlite3
    sudo apt-get install jq
    sudo apt-get install zip
    sudo apt-get install build-essential

    // install stackdriver
    curl -sSO https://dl.google.com/cloudagents/add-monitoring-agent-repo.sh && sudo bash add-monitoring-agent-repo.sh && sudo apt-get update && sudo apt-cache madison stackdriver-agent

    sudo apt-get install -y 'stackdriver-agent=6.*'

    // version check
    dpkg-query --show --showformat '${Package} ${Version} ${Architecture} ${Status}\n' stackdriver-agent

## Setup MySql
    mysql_secure_installation
    // run this query to require root password when connecting to db
    USE mysql; 
    UPDATE mysql.user SET plugin = 'mysql_native_password' WHERE user = 'root' AND host = 'localhost'; 
    UPDATE user set authentication_string=PASSWORD("mynewpassword") where User='root'; 
    FLUSH PRIVILEGES;
    exit;
    sudo service mysql restart



## Install Go
    cd /tmp
    wget https://dl.google.com/go/go1.16.3.linux-amd64.tar.gz

    sudo tar -xvf go1.16.3.linux-amd64.tar.gz
    sudo mv go /usr/local

    sudo nano ~/.bashrc
    // add to end of file
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

    source ~/.profile

    go version

## Clone App from Git
    cd $GOPATH
    mkdir src && cd src
    mkdir github.com && cd github.com
    mkdir btcid && cd btcid

    git clone git@35.240.159.3:btcid/wallet/wallet-services-backend-go.git
    cd wallet-services-backend-go
    mkdir logs

## Run HTTP App
    go mod init

    // build app
    go build cmd/server/main.go

    // setup database

    // run app
    ./main

## Run Cron App
    // build app
    go build cmd/server/main.go

    // setup database

    // run app
    // sleep duration e.g. "300ms", "1.5m" or "60s". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h" 
    ./main -app cron -function [[function name]] -sleep [[ sleep duration ]]
