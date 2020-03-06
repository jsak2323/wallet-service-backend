// INSTALL VM REQUIREMENTS
    export LANGUAGE=en_US.UTF-8
    export LC_ALL=en_US.UTF-8

    sudo apt-get update
    sudo apt-get upgrade

    sudo apt-get install apache2
    sudo apt-get install mysql-server 
    //user & password: root 123456

    sudo apt-get update && sudo apt-get install sqlite3

    sudo a2enmod rewrite

    // Change AllowOverride None to AllowOverride All in apache2.conf

    sudo service apache2 restart

    // install go
    cd /tmp
    wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    sudo tar -xvf go1.13.linux-amd64.tar.gz
    sudo mv go /usr/local

    // add to ~/.profile
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

    source ~/.profile


// INSTALL APP
    cd
    mkdir go && cd go && mkdir src && cd src && mkdir github.com && cd github.com && mkdir btcid && cd btcid

    // clone git repository

    cd wallet-services-backend
    go mod init

    // build app
    go build cmd/server/main.go

    // run app
    ./main







