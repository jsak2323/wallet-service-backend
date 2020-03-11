CREATE DATABASE wallet_service;
USE wallet_service;

CREATE TABLE currency_config (
  id                        int(11) NOT NULL AUTO_INCREMENT,
  symbol                    varchar(10) NOT NULL,
  name                      varchar(50) NOT NULL,
  name_uppercase            varchar(50) NOT NULL,
  name_lowercase            varchar(50) NOT NULL,
  unit                      varchar(10) NOT NULL,
  token_type                varchar(10) NOT NULL DEFAULT 'main',
  is_finance_enabled        tinyint(1) NOT NULL DEFAULT 0,
  is_single_address         tinyint(1) NOT NULL DEFAULT 0,
  is_using_memo             tinyint(1) NOT NULL DEFAULT 0,
  is_qrcode_enabled         tinyint(1) NOT NULL DEFAULT 0,
  is_address_notice_enabled tinyint(1) NOT NULL DEFAULT 0,
  qrcode_prefix             varchar(50) NULL DEFAULT NULL,
  withdraw_fee              varchar(50) NOT NULL DEFAULT "0",
  default_idr_price         int(15) NOT NULL DEFAULT 0,
  cmc_id                    int(7) NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_config (
  id            int(11) NOT NULL AUTO_INCREMENT,
  currency_id   int(11) NOT NULL,
  type          varchar(30) NOT NULL,
  host          varchar(30) NOT NULL,
  port          varchar(30) NOT NULL,
  `path`        varchar(30) NOT NULL,
  user          varchar(30) NOT NULL,
  password      varchar(100) NOT NULL,
  hashkey       varchar(150) NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (currency_id) REFERENCES currency_config(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;