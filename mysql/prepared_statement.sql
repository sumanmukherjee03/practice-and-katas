--  ##############################################################################################
--  ######################################## PREPARED STATEMENT ##################################
--  ##############################################################################################

PREPARE stmt FROM
  'SELECT productCode, productName FROM products WHERE productCode = ?';
SET @pc = 'S10_1678';
EXECUTE stmt USING @pc;

DEALLOCATE PREPEARE stmt;
