--  ##############################################################################################
--  ############################################ TRIGGERS ########################################
--  ##############################################################################################
CREATE TABLE students (st_roll INT, age INT, name VARCHAR(30), marks FLOAT);

SHOW TRIGGERS FROM testdb LIKE '%set_default%';

--  Example of creating a before insert trigger
DELIMITER $$
CREATE TRIGGER set_default_marks
BEFORE INSERT ON students
FOR EACH ROW
BEGIN
  IF new.marks < 0 THEN
    SET new.marks = 40;
  END IF;
END$$
DELIMITER ;

--  The insert statement will call the trigger
INSERT INTO students VALUES(501, 10, 'Jane', 55.0),
  (502, 10, 'John', -1);

--  To delete a trigger
DROP TRIGGER set_default_marks;


--  In an insert trigger you can access the NEW modifiers
DELIMITER $$
CREATE TRIGGER after_member_inserts
AFTER INSERT
ON members
FOR EACH ROW
BEGIN
  IF NEW.birthDate IS NULL THEN
    INSERT INTO reminders(memberId, message)
    VALUES(NEW.id, CONCAT('Hi', NEW.name, ', please update your birth date'));
  END IF;
END$$
DELIMITER ;


--  In a update trigger you can access both the NEW and OLD modifiers
--  Also in a before update trigger you can update the new values but you cant update the old values
DELIMITER $$
CREATE TRIGGER before_sales_update
BEFORE UPDATE
ON sales FOR EACH ROW
BEGIN
    DECLARE errorMessage VARCHAR(255);
    SET errorMessage = CONCAT('The new quantity ', NEW.quantity, ' cannot be 3 times greater than the current quantity ', OLD.quantity);
    IF NEW.quantity > OLD.quantity * 3 THEN
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = errorMessage;
    END IF;
END $$
DELIMITER ;

--  This is a simple example of an after update trigger.
DELIMITER $$
CREATE TRIGGER after_sales_update
AFTER UPDATE
ON sales FOR EACH ROW
BEGIN
    IF OLD.quantity <> NEW.quantity THEN
      INSERT INTO SalesChanges(salesId, beforeQuantity, afterQuantity) VALUES(OLD.id, OLD.quantity, NEW.quantity);
    END IF;
END$$
DELIMITER ;



--  In a before delete trigger you can access the old values via the OLD modifier
--  if you want to populate some message or archive things
--  An after delete trigger also has access to the OLD modifier
DELIMITER $$
CREATE TRIGGER before_salaries_delete
BEFORE DELETE
ON salaries FOR EACH ROW
BEGIN
    INSERT INTO SalaryArchives(employeeNumber,validFrom,amount) VALUES(OLD.employeeNumber,OLD.validFrom,OLD.amount);
END$$
DELIMITER ;
