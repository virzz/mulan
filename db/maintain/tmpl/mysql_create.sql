CREATE USER IF NOT EXISTS '{{.User}}'@'%' identified by '{{.Pass}}';
CREATE DATABASE IF NOT EXISTS `{{.Name}}` CHARACTER SET = 'utf8mb4' COLLATE = 'utf8mb4_unicode_ci';
GRANT ALL ON `{{.Name}}`.* to '{{.User}}'@'%';
FLUSH PRIVILEGES;