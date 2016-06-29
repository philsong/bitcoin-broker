# ************************************************************
# Sequel Pro SQL dump
# Version 4135
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.5.38)
# Database: trader
# Generation Time: 2015-07-31 06:14:27 +0000
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
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL,
  `available_cny` decimal(65,2) NOT NULL,
  `available_btc` decimal(65,4) NOT NULL,
  `frozen_cny` decimal(65,2) NOT NULL,
  `frozen_btc` decimal(65,4) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `exchange` (`exchange`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table depth
# ------------------------------------------------------------

CREATE TABLE `depth` (
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL,
  `orderbook` varchar(4096) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table exchange_config
# ------------------------------------------------------------

CREATE TABLE `exchange_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL,
  `accessKey` varchar(64) NOT NULL,
  `secretKey` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `exchange` (`exchange`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table site_order
# ------------------------------------------------------------

CREATE TABLE `site_order` (
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `client_id` varchar(64) NOT NULL,
  `parent_id` int(11) NOT NULL,
  `trade_type` varchar(32) NOT NULL,
  `order_status` varchar(32) NOT NULL,
  `amount` decimal(65,4) NOT NULL,
  `estimate_price` decimal(65,2) NOT NULL,
  `estimate_cny` decimal(65,2) NOT NULL,
  `estimate_btc` decimal(65,4) NOT NULL,
  `info` longtext NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id` (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table ticker
# ------------------------------------------------------------

CREATE TABLE `ticker` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ask` decimal(65,2) NOT NULL,
  `bid` decimal(65,2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table trade_order
# ------------------------------------------------------------

CREATE TABLE `trade_order` (
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exchange` varchar(64) NOT NULL,
  `trade_type` varchar(32) NOT NULL,
  `order_status` varchar(32) NOT NULL,
  `estimate_cny` decimal(65,2) NOT NULL,
  `estimate_btc` decimal(65,4) NOT NULL,
  `estimate_price` decimal(65,2) NOT NULL,
  `deal_cny` decimal(65,2) NOT NULL,
  `deal_btc` decimal(65,4) NOT NULL,
  `deal_price` decimal(65,2) NOT NULL,
  `order_id` varchar(64) NOT NULL,
  `info` longtext NOT NULL,
  `parent_id` int(11) NOT NULL,
  `site_order_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `trade_order_parent_id_665aa13862248bbc_fk_site_order_id` (`parent_id`),
  KEY `trade_order_site_order_id_532cfc808c8633fc_fk_site_order_id` (`site_order_id`),
  CONSTRAINT `trade_order_site_order_id_532cfc808c8633fc_fk_site_order_id` FOREIGN KEY (`site_order_id`) REFERENCES `site_order` (`id`),
  CONSTRAINT `trade_order_parent_id_665aa13862248bbc_fk_site_order_id` FOREIGN KEY (`parent_id`) REFERENCES `site_order` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
