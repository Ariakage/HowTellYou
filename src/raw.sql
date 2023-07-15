-- Create User Table
CREATE TABLE
IF
	NOT EXISTS hty_user (
		`id` INT PRIMARY KEY AUTO_INCREMENT,
		`favimg` TEXT NOT NULL,
		`name` VARCHAR ( 16 ) NOT NULL,
		`nickname` VARCHAR ( 20 ) NOT NULL,
		`email` VARCHAR ( 50 ) NOT NULL,
		`pwd` VARCHAR ( 512 ) NOT NULL,
	`create_time` DATETIME DEFAULT CURRENT_TIMESTAMP 
	);
-- Create Friend Table (https://blog.csdn.net/wo541075754/article/details/82733278)
CREATE TABLE
IF
	NOT EXISTS hty_friend (
		`user_id` INT NOT NULL,
		`friend_id` INT NOT NULL,
		`user_group` VARCHAR ( 10 ) NOT NULL,
	`friend_group` VARCHAR ( 10 ) NOT NULL 
	);
-- Create Group Table (https://blog.csdn.net/php_xml/article/details/108690219)
CREATE TABLE
IF
	NOT EXISTS hty_group (
		`id` INT PRIMARY KEY AUTO_INCREMENT,
		`favimg` TEXT DEFAULT '',
		`name` VARCHAR ( 16 ) NOT NULL,
		`owner_id` INT NOT NULL,
		`admins` LONGTEXT NOT NULL DEFAULT '',
		`members` LONGTEXT NOT NULL,
		`type` INT NOT NULL,
		`remark` VARCHAR ( 200 ) NOT NULL DEFAULT '',
	`create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL 
	);
-- Create Message Table (https://blog.csdn.net/qq_42249896/article/details/104033697)
CREATE TABLE
IF
	NOT EXISTS hty_message (
		`id` INT PRIMARY KEY AUTO_INCREMENT,
		`send_user_id` INT NOT NULL,
		`receive_user_id` INT NOT NULL,
		`content` TEXT NOT NULL,
	`send_time` DATETIME NOT NULL 
	);

-- User

-- Add User
INSERT INTO hty_user(`favimg` ,`name`, `nickname`, `email`, `pwd`) VALUES ('', 'test_user1', 'test_user1','abcd@test.com', '114514')
-- Select User pwd
SELECT `pwd` FROM hty_user WHERE `id` = 11 or `email` = 'abcd@test.com'

-- End User