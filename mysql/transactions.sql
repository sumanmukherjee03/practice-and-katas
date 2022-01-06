--  ##############################################################################################
--  ########################################## TRANSACTIONS ######################################
--  ##############################################################################################

--  To automatically commit changes turn on autocommit
SET autocommit = 0;

--  To stop the database from automatically committing changes turn off autocommit
SET autocommit = 1;


--  The start of a transaction can be denoted by `START TRANSACTION` or `BEGIN` or `BEGIN WORK`
--  Here's a script that you can probably put in a function to run a transaction
--  Sart of the transaction script
START TRANSACTION;

DECLARE nextOrderNumber INT DEFAULT 0
SELECT
  MAX(orderNumber)+1 INTO nextOrderNumber
  FROM orders;

INSERT INTO orders(orderNumber, orderDate, requiredDate, shippedDate, status, customerNumber)
  VALUES(nextOrderNumber, '2022-01-03', '2022-01-05', '2022-01-06', 'In Process', 150);

INSERT INTO orderdetails(orderNumber, productCode, quantityOrdered, priceEach, orderLineNumber)
  VALUES(nextOrderNumber, 'SMPY54T21', 30, '128', 1),
        (nextOrderNumber, 'MTNL33Q87', 10, '772', 3);

COMMIT;
--  End of the transaction script



--  To get the newly created sales order from the previously created transaction
SELECT a.orderName, orderDate, requiredDate, shippedDate, status, comments, customerNumber, orderLineNumber, productCode, quantityOrdered, priceEach
  FROM orders a
  INNER JOIN orderdetails b USING (orderNumber)
  WHERE a.orderNumber = 445;


--  ############################################################################
--  Sample case of using a transaction in a stored proc to automatically insert rows into join tables with foreign keys
--  ############################################################################
CREATE TABLE accounts (
  account_id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL
);

CREATE TABLE phones (
  phone_id INT AUTO_INCREMENT,
  account_id INT NOT NULL,
  phone VARCHAR(25) NOT NULL,
  description VARCHAR(255) NOT NULL,
  PRIMARY KEY (phone_id, account_id),
  FOREIGN KEY (account_id) REFERENCES accounts (account_id)
);

DELIMITER $$
CREATE PROCEDURE createAccount(fname VARCHAR(255), lname VARCHAR(255), phone VARCHAR(25), description VARCHAR(255))
BEGIN
  DECLARE l_account_id INT  DEFAULT 0;
  START TRANSACTION;
  INSERT INTO accounts(first_name, last_name) VALUES(fname, lname);
  SET l_account_id = LAST_INSERT_ID();
  IF l_account_id > 0
    INSERT INTO phones(account_id, phone, description) VALUES(l_account_id, phone, description);
    COMMIT;
  ELSE
    ROLLBACK;
  END IF;
END$$
DELIMITER ;
