--  ##############################################################################################
--  ############################################ INDEXES #########################################
--  ##############################################################################################

CREATE INDEX idxJobTitle ON employees(jobTitle);

SHOW INDEXES FROM employees IN testdb;
SHOW INDEXES FROM employees WHERE VISIBLE = 'NO';

--  when you run the query below with the index in place for jobTitle you will see less rows scanned in the explain output.
--  Also, it would indicate which possible indexes could be used and which indexes were actually used - possible_keys and keys column in the output of explain.
EXPLAIN SELECT employeeNumber, lastName, firstName FROM employees WHERE jobTitle = 'Sales Rep';

--  The types of algorithm available are COPY and INPLACE.
  --  COPY - The table is copied to the new table row by row, the DROP INDEX is then performed on the copy of the original table.
  --    The concurrent data manipulation statements such as INSERT and UPDATE are not permitted.
  --  INPLACE - The table is rebuilt in place instead of copied to the new one.
  --    MySQL issues an exclusive metadata lock on the table during the preparation and execution phases of the index removal operation. This algorithm allows for concurrent data manipulation statements.
--  There are also multiple types of locks - DEFAULT, EXCLUSIVE, SHARED and NONE
DROP INDEX idxJobTitle ON employees ALGORITHM = INPLACE LOCK = DEFAULT;

CREATE INDEX idx_extension ON employees(extension) INVISIBLE;
ALTER TABLE employees ALTER INDEX idx_extension VISIBLE;

--  If you have multiple indexes on a few columns and want certain queries to use some specific indexes
--  then specify that using the 'USE INDEX' clause. If you put an EXPLAIN infront of the query you can see that in the possible_keys and keys columns.
SELECT * FROM customers USE INDEX (idx_name_fl, idx_name_lf) WHERE contactFirstName LIKE 'A%' OR contactLastName LIKE 'A%';
