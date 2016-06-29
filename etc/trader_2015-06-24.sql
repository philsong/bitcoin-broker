# ************************************************************
# Sequel Pro SQL dump
# Version 4135
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.5.38)
# Database: trader
# Generation Time: 2015-06-24 04:54:41 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table account
# ------------------------------------------------------------

CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL DEFAULT '',
  `available_cny` double(65,2) NOT NULL DEFAULT '0.00',
  `available_btc` double(65,4) NOT NULL DEFAULT '0.0000',
  `frozen_cny` double(65,2) NOT NULL DEFAULT '0.00',
  `frozen_btc` double(65,4) NOT NULL DEFAULT '0.0000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `exchange` (`exchange`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table depth
# ------------------------------------------------------------

CREATE TABLE `depth` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL DEFAULT '',
  `orderbook` varchar(4096) NOT NULL DEFAULT '',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table exchange_config
# ------------------------------------------------------------

CREATE TABLE `exchange_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL DEFAULT '',
  `accessKey` varchar(64) NOT NULL DEFAULT '',
  `secretKey` varchar(64) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `exchange` (`exchange`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table site_order
# ------------------------------------------------------------

CREATE TABLE `site_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `client_id` varchar(64) NOT NULL DEFAULT '',
  `parent_id` bigint(20) NOT NULL DEFAULT '0',
  `trade_type` varchar(32) NOT NULL DEFAULT '',
  `order_status` varchar(32) NOT NULL DEFAULT '',
  `amount` double(65,4) NOT NULL DEFAULT '0.0000',
  `estimate_price` double(65,2) NOT NULL DEFAULT '0.00',
  `estimate_cny` double(65,2) NOT NULL DEFAULT '0.00',
  `estimate_btc` double(65,4) NOT NULL DEFAULT '0.0000',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `info` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table ticker
# ------------------------------------------------------------

CREATE TABLE `ticker` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ask` double(65,2) NOT NULL DEFAULT '0.00',
  `bid` double(65,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table trade_order
# ------------------------------------------------------------

CREATE TABLE `trade_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `site_order_id` bigint(20) NOT NULL DEFAULT '0',
  `parent_id` bigint(20) NOT NULL DEFAULT '0',
  `exchange` varchar(64) NOT NULL DEFAULT '',
  `trade_type` varchar(32) NOT NULL DEFAULT '',
  `order_status` varchar(32) NOT NULL DEFAULT '',
  `estimate_cny` double(65,2) NOT NULL DEFAULT '0.00',
  `estimate_btc` double(65,4) NOT NULL DEFAULT '0.0000',
  `estimate_price` double(65,2) NOT NULL DEFAULT '0.00',
  `deal_cny` double(65,2) NOT NULL DEFAULT '0.00',
  `deal_btc` double(65,4) NOT NULL DEFAULT '0.0000',
  `deal_price` double(65,2) NOT NULL DEFAULT '0.00',
  `order_id` varchar(64) NOT NULL DEFAULT '0',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `info` text,
  PRIMARY KEY (`id`),
  KEY `site_order_id` (`site_order_id`),
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `trade_order_ibfk_1` FOREIGN KEY (`site_order_id`) REFERENCES `site_order` (`id`),
  CONSTRAINT `trade_order_ibfk_2` FOREIGN KEY (`parent_id`) REFERENCES `site_order` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
