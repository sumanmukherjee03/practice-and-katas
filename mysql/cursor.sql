--  ##############################################################################################
--  ############################################ CURSORS #########################################
--  ##############################################################################################

-- Cursons allow you to iterate over result sets of a query, for example results from a stored proc
--  Mysql cursors are read-only, non-scrollable and asentitive. Non-scrollable means that we cant jump to
--  a particular row or reverse order etc. Order is determined by the select clause. Asensitive means that
--  the rows of the result are actual rows in table and not a copy. So, if any other mysql conn updates the rows
--  that would be visible in the result set.

DELIMITER $$
CREATE PROCEDURE createEmailList (INOUT emailList VARCHAR(5000))
BEGIN
  --  In this example below, forst we declare 2 variables.
  --  Cursor declaration must be after variable declaration.
  DECLARE finished INT DEFAULT 0;
  DECLARE emailAddr VARCHAR(100) DEFAULT "";
  --  Declare the cursor for the select statement
  DECLARE cursorEmail CURSOR FOR SELECT email FROM employees;
  --  declare NOT FOUND handler
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET finished = 1;
  --  open cursor
  OPEN cursorEmail;
  --  Iterate over each email in a loop by fetching based on the cursor position
  --  and create a list of concatenated emails separated by semicolon.
  getEmail: LOOP
    FETCH cursorEmail INTO emailAddr;
    IF finished = 1 THEN
      LEAVE getEmail;
    END IF;
    SET emailList = CONCAT(emailAddr, ";", emailLis);
  END LOOP getEmail;
  --  close cursor
  CLOSE cursorEmail;
END$$
DELIMITER ;
