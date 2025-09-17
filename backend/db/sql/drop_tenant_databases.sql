-- Drops all databases starting with 'ednevnik_' except 'ednevnik_workspace'
-- Run this as a user with privileges to drop databases

DELIMITER $$

CREATE PROCEDURE drop_ednevnik_tenant_dbs()
BEGIN
  DECLARE done INT DEFAULT FALSE;
  DECLARE dbname VARCHAR(255);
  DECLARE cur CURSOR FOR
    SELECT SCHEMA_NAME
    FROM INFORMATION_SCHEMA.SCHEMATA
    WHERE SCHEMA_NAME LIKE 'ednevnik\_%'
      AND SCHEMA_NAME != 'ednevnik_workspace';
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

  OPEN cur;
  read_loop: LOOP
    FETCH cur INTO dbname;
    IF done THEN
      LEAVE read_loop;
    END IF;
    -- Log prior to dropping database
    SELECT CONCAT('Dropping database: ', dbname) AS log_message;

    SET @s = CONCAT('DROP DATABASE IF EXISTS `', dbname, '`');
    PREPARE stmt FROM @s;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- Log after dropping database
    SELECT CONCAT('Database dropped: ', dbname) AS log_message;
  END LOOP;
  CLOSE cur;
END$$

DELIMITER ;

CALL drop_ednevnik_tenant_dbs();

DROP PROCEDURE drop_ednevnik_tenant_dbs;

SHOW DATABASES LIKE 'ednevnik\_%';
