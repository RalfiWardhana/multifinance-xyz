-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: kredit_plus
-- ------------------------------------------------------
-- Server version	8.0.42-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `customer_limits`
--

DROP TABLE IF EXISTS `customer_limits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `customer_limits` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `customer_id` bigint NOT NULL,
  `tenor_months` bigint NOT NULL,
  `limit_amount` decimal(15,2) NOT NULL,
  `used_amount` double DEFAULT '0',
  `available_amount` decimal(15,2) GENERATED ALWAYS AS ((`limit_amount` - `used_amount`)) STORED,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_customer_tenor` (`customer_id`,`tenor_months`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_tenor_months` (`tenor_months`),
  KEY `idx_customer_limits_customer_id` (`customer_id`),
  KEY `idx_customer_limits_tenor_months` (`tenor_months`),
  CONSTRAINT `customer_limits_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_customer_limits_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customer_limits`
--

LOCK TABLES `customer_limits` WRITE;
/*!40000 ALTER TABLE `customer_limits` DISABLE KEYS */;
INSERT INTO `customer_limits` (`id`, `customer_id`, `tenor_months`, `limit_amount`, `used_amount`, `created_at`, `updated_at`) VALUES (25,10,1,1000000.00,0,'2025-07-14 10:06:33.740','2025-07-14 10:06:33.740'),(26,10,2,2000000.00,0,'2025-07-14 10:06:33.746','2025-07-14 10:06:33.746'),(27,10,3,3000000.00,0,'2025-07-14 10:06:33.750','2025-07-14 10:06:33.750'),(28,10,4,4000000.00,0,'2025-07-14 10:06:33.757','2025-07-14 10:06:33.757');
/*!40000 ALTER TABLE `customer_limits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS `customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `customers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `nik` varchar(16) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `legal_name` varchar(255) NOT NULL,
  `birth_place` varchar(255) NOT NULL,
  `birth_date` date NOT NULL,
  `salary` decimal(15,2) NOT NULL,
  `ktp_photo_path` varchar(500) DEFAULT NULL,
  `selfie_photo_path` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_customers_nik` (`nik`),
  UNIQUE KEY `uni_customers_user_id` (`user_id`),
  KEY `idx_nik` (`nik`),
  KEY `idx_full_name` (`full_name`),
  KEY `idx_customers_deleted_at` (`deleted_at`),
  KEY `idx_customers_nik` (`nik`),
  KEY `idx_customers_user_id` (`user_id`),
  CONSTRAINT `fk_customers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customers`
--

LOCK TABLES `customers` WRITE;
/*!40000 ALTER TABLE `customers` DISABLE KEYS */;
INSERT INTO `customers` VALUES (10,'3174012345678901','Ralfi Wardhana','Ralfi Wardhana','Jakarta','1996-01-01',9000000.00,'/uploads/ktp/ralfi_ktp_20250714.jpg','/uploads/selfie/ralfi_selfie_20250714.jpg','2025-07-14 10:06:33.735','2025-07-14 10:06:33.735',NULL,4);
/*!40000 ALTER TABLE `customers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `contract_number` varchar(50) NOT NULL,
  `customer_id` bigint NOT NULL,
  `tenor_months` bigint NOT NULL,
  `otr_amount` decimal(15,2) NOT NULL,
  `admin_fee` decimal(15,2) NOT NULL,
  `installment_amount` decimal(15,2) NOT NULL,
  `interest_amount` decimal(15,2) NOT NULL,
  `asset_name` varchar(255) NOT NULL,
  `asset_type` enum('WHITE_GOODS','MOTOR','MOBIL') NOT NULL,
  `status` enum('PENDING','APPROVED','REJECTED','ACTIVE','COMPLETED','DEFAULTED') DEFAULT 'PENDING',
  `transaction_source` enum('ECOMMERCE','WEB','DEALER') NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_transactions_contract_number` (`contract_number`),
  KEY `idx_contract_number` (`contract_number`),
  KEY `idx_customer_id` (`customer_id`),
  KEY `idx_status` (`status`),
  KEY `idx_asset_type` (`asset_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_transactions_deleted_at` (`deleted_at`),
  KEY `idx_transactions_contract_number` (`contract_number`),
  KEY `idx_transactions_customer_id` (`customer_id`),
  KEY `idx_transactions_asset_type` (`asset_type`),
  KEY `idx_transactions_status` (`status`),
  KEY `idx_transactions_created_at` (`created_at`),
  CONSTRAINT `fk_transactions_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`),
  CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `role` enum('ADMIN','CUSTOMER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'CUSTOMER',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`),
  UNIQUE KEY `uni_users_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_username` (`username`),
  KEY `idx_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (4,'ralfi_customer','ralfi@example.com','$2a$10$S9UI0ybjgzScHUk/s3jWAemHeEpB8r33YZWLsk8ny8xdwrRVwac6S','CUSTOMER',1,'2025-07-14 10:06:33.726','2025-07-14 10:06:33.726',NULL),(5,'admin','admin@ptxyz.com','$2a$10$K5QzJ8gOLJ8K5QzJ8gOLJO5QzJ8gOLJ8K5QzJ8gOLJ8K5QzJ8gOLJO','ADMIN',1,'2025-07-14 10:43:50.220','2025-07-14 10:43:50.220',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-07-14 11:15:26
