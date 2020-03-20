CREATE DATABASE wallet_service;
USE wallet_service;

CREATE TABLE currency_config (
  id                        INT(11) NOT NULL AUTO_INCREMENT,
  symbol                    VARCHAR(10) NOT NULL,
  name                      VARCHAR(50) NOT NULL,
  name_uppercase            VARCHAR(50) NOT NULL,
  name_lowercase            VARCHAR(50) NOT NULL,
  unit                      VARCHAR(10) NOT NULL,
  token_type                VARCHAR(10) NOT NULL DEFAULT 'main',
  is_finance_enabled        TINYINT(1) NOT NULL DEFAULT 0,
  is_single_address         TINYINT(1) NOT NULL DEFAULT 0,
  is_using_memo             TINYINT(1) NOT NULL DEFAULT 0,
  is_qrcode_enabled         TINYINT(1) NOT NULL DEFAULT 0,
  is_address_notice_enabled TINYINT(1) NOT NULL DEFAULT 0,
  qrcode_prefix             VARCHAR(50) NULL DEFAULT NULL,
  withdraw_fee              VARCHAR(50) NOT NULL DEFAULT "0",
  healthy_block_diff        INT(11) NOT NULL DEFAULT 0,
  default_idr_price         INT(15) NOT NULL DEFAULT 0,
  cmc_id                    INT(7) NULL DEFAULT NULL,

  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_config (
  id                        INT(11) NOT NULL AUTO_INCREMENT,
  currency_id               INT(11) NOT NULL,
  type                      VARCHAR(30) NOT NULL,
  host                      VARCHAR(30) NOT NULL,
  port                      VARCHAR(30) NOT NULL,
  `path`                    VARCHAR(30) NOT NULL,
  user                      VARCHAR(30) NOT NULL,
  password                  VARCHAR(100) NOT NULL,
  hashkey                   VARCHAR(150) NOT NULL,
  node_version              VARCHAR(50) NOT NULL,
  node_last_updated         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  is_health_check_enabled   TINYINT(1) NOT NULL DEFAULT 0,

  PRIMARY KEY (id),
  FOREIGN KEY (currency_id) REFERENCES currency_config(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE health_check (
  id                    INT(11) NOT NULL AUTO_INCREMENT,
  rpc_config_id         INT(11) NOT NULL,
  blockcount            INT(11) NOT NULL DEFAULT 0,
  block_diff            INT(11) NOT NULL DEFAULT 0,
  is_healthy            TINYINT(1) NOT NULL DEFAULT 1,
  last_updated          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

  PRIMARY KEY (id),
  FOREIGN KEY (rpc_config_id) REFERENCES rpc_config(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

