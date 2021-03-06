
USE wallet_service;

CREATE TABLE currency_config (
  id                        INT(11) NOT NULL AUTO_INCREMENT,
  symbol                    VARCHAR(10) NOT NULL,
  name                      VARCHAR(50) NOT NULL,
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
  parent_symbol             VARCHAR(50) NULL DEFAULT NULL,
  address                   VARCHAR(255) NOT NULL DEFAULT "",
  module_type               VARCHAR(50) NOT NULL DEFAULT "",
  active                    TINYINT(1) NOT NULL DEFAULT 0,
  last_updated              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

  PRIMARY KEY (id),
  CONSTRAINT symbol_token_type UNIQUE (symbol,token_type)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_config (
  id                        INT(11) NOT NULL AUTO_INCREMENT,
  type                      VARCHAR(30) NOT NULL,
  name                      VARCHAR(50) NOT NULL DEFAULT "",
  platform                  VARCHAR(30) NOT NULL DEFAULT "GCP",
  host                      VARCHAR(30) NOT NULL,
  port                      VARCHAR(30) NOT NULL,
  `path`                    VARCHAR(30) NOT NULL,
  user                      VARCHAR(30) NOT NULL,
  password                  VARCHAR(100) NOT NULL,
  hashkey                   VARCHAR(150) NOT NULL,
  node_version              VARCHAR(50) NOT NULL,
  node_last_updated         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  is_health_check_enabled   TINYINT(1) NOT NULL DEFAULT 0,
  atom_feed                 VARCHAR(255) NOT NULL DEFAULT "",
  active                    TINYINT(1) NOT NULL DEFAULT 0,

  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE currency_config_rpc_config (
    currency_config_id INT(11) NOT NULL,
    rpc_config_id INT(11) NOT NULL,

    FOREIGN KEY (currency_config_id) REFERENCES currency_config(id),
    FOREIGN KEY (rpc_config_id) REFERENCES rpc_config(id),
    CONSTRAINT currency_rpc UNIQUE (currency_config_id, rpc_config_id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE system_config (
  name    VARCHAR(255) NOT NULL DEFAULT "",
  value   TEXT NOT NULL,

  UNIQUE KEY name (name)
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

CREATE TABLE cold_balance (
  id                    INT(11) NOT NULL AUTO_INCREMENT,
  currency_id           INT(11) NOT NULL,
  name                  VARCHAR(50) NOT NULL DEFAULT "",
  type                  VARCHAR(20) NOT NULL DEFAULT "",
  fireblocks_name       VARCHAR(50) NOT NULL DEFAULT "",
  balance               VARCHAR(255) NOT NULL DEFAULT "",
  address               VARCHAR(255) NOT NULL DEFAULT "",
  active                TINYINT(1) NOT NULL DEFAULT 0,
  last_updated          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

  PRIMARY KEY (id),
  UNIQUE KEY name (name),
  FOREIGN KEY (currency_id) REFERENCES currency_config(id)
 ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE hot_limit (
   id           INT(11) NOT NULL AUTO_INCREMENT,
   currency_id  INT(11) NOT NULL,
   type         VARCHAR(20) NOT NULL DEFAULT "",
   amount       VARCHAR(255) NOT NULL DEFAULT "",

   PRIMARY KEY (id),
   FOREIGN KEY (currency_id) REFERENCES currency_config(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_method (
    id          INT(11) NOT NULL AUTO_INCREMENT,
    name        VARCHAR(50) NOT NULL DEFAULT "",
    type        VARCHAR(50) NOT NULL DEFAULT "",
    network     VARCHAR(50) NOT NULL DEFAULT "",
    num_of_args INT(11) NOT NULL DEFAULT 6,

    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_config_rpc_method (
    rpc_config_id INT(11) NOT NULL,
    rpc_method_id INT(11) NOT NULL,

    FOREIGN KEY (rpc_config_id) REFERENCES rpc_config(id),
    FOREIGN KEY (rpc_method_id) REFERENCES rpc_method(id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE rpc_request (
    id                INT(11) NOT NULL AUTO_INCREMENT,
    arg_name          VARCHAR(50) NOT NULL DEFAULT "",
    type              VARCHAR(50) NOT NULL DEFAULT "",
    arg_order         INT(11) NOT NULL DEFAULT -1,
    source            VARCHAR(50) NOT NULL DEFAULT "",
    runtime_var_name  VARCHAR(50) NOT NULL DEFAULT "",
    value             VARCHAR(255) NOT NULL DEFAULT "",
    rpc_method_id     INT(11) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (rpc_method_id) REFERENCES rpc_method(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE rpc_response (
    id            INT(11) NOT NULL AUTO_INCREMENT,
    xml_path      VARCHAR(255) NOT NULL DEFAULT "",
    field_name    VARCHAR(50) NOT NULL DEFAULT "",
    data_type_tag VARCHAR(50) NOT NULL DEFAULT "",
    parse_type    VARCHAR(50) NOT NULL DEFAULT "",
    json_fields   VARCHAR(255) NOT NULL DEFAULT "",
    rpc_method_id INT(11) NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (rpc_method_id) REFERENCES rpc_method(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE deposit (
  id            BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  currency_id   INT(11) NOT NULL,
  address_to    VARCHAR(255) NOT NULL,
  tx            VARCHAR(255) NOT NULL,
  memo          VARCHAR(255) NOT NULL,
  log_index     VARCHAR(255) NOT NULL,
  confirmations INT(11) NOT NULL DEFAULT 0,
  amount        VARCHAR(255) NOT NULL DEFAULT "",
  success_time  DATETIME NULL,
  last_updated  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

  PRIMARY KEY (id),
  FOREIGN KEY (currency_id) REFERENCES currency_config(id),
  UNIQUE KEY (tx),
  INDEX (currency_id, address_to),
  INDEX (currency_id, tx),
  INDEX (currency_id, success_time)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE withdraw (
  id             BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  currency_id    INT(11) NOT NULL,
  address_to     VARCHAR(255) NOT NULL,
  tx             VARCHAR(255) NOT NULL,
  memo           VARCHAR(255) NOT NULL,
  confirmations  INT(11) NOT NULL,
  blockchain_fee VARCHAR(255) NOT NULL,
  market_price   VARCHAR(255) NOT NULL,
  log_index      VARCHAR(255) NOT NULL,
  amount         VARCHAR(255) NOT NULL,
  success_time   DATETIME NULL,
  last_updated   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),

  PRIMARY KEY (id),
  FOREIGN KEY (currency_id) REFERENCES currency_config(id),
  UNIQUE KEY (tx),
  INDEX (currency_id, address_to),
  INDEX (currency_id, tx),
  INDEX (currency_id, success_time)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;