--  ##############################################################################################
--  ############################################ FUNCTIONS #######################################
--  ##############################################################################################

--  functions have to be deterministic or non-deterministic. If the same set of inputs always produce the
--  same output, then it is deterministic. More often than not mysql functions are non-deterministic in nature.
--  So, thats the default in mysql function types.
--  A key difference between functions and stored procs is that a function can be invoked like an expression SUM(columnName).
--  it doesnt require the CALL keyword like a stored proc.

SELECT routine_name FROM information_schema.routines WHERE routine_type = 'FUNCTION' AND routine_schema = '<database_name>';
SHOW FUNCTION STATUS LIKE '%CustomerLevel%';

DELIMITER $$
CREATE FUNCTION CustomerLevel(
  credit DECIMAL(10,2)
)
RETURNS VARCHAR(20)
DETERMINISTIC
BEGIN
  DECLARE custLevel VARCHAR(20);
  IF credit > 50000 THEN
    SET custLevel = 'PLATINUM';
  ELSEIF (credit >= 10000 AND credit <= 50000) THEN
    SET custLevel = 'GOLD';
  ELSEIF credit < 10000 THEN
    SET custLevel = 'SILVER';
  END IF;
  RETURN custLevel;
END$$
DELIMITER ;

SELECT customerName, CustomerLevel(creditLimit) FROM customers ORDER BY customerName;


DROP FUNCTION IF EXISTS CustomerLevel;
