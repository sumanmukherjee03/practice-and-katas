--  ##############################################################################################
--  ########################################## TABLE LOCKS #######################################
--  ##############################################################################################

--  ######################################## READ LOCK ##################################
--  This returns the connection id of the current connection
SELECT CONNECTION_ID();

--  Once a READ lock is acquired on a table you can not write to that table within the same mysql session. It will error out immediately.
--  However, you can attempt to write from a different session in which case it will wait until the READ lock is released by the first session.
--  If you query the processlist, you will see the insert statement as "Waiting for table metadata lock"
--  However other sessions can read data from the table even when there is a READ lock on the table.
LOCK TABLES messages READ;

--  Once you release the READ lock on the table, the insert statement will go through
UNLOCK TABLES;


--  ######################################## WRITE LOCK ##################################
--  Once a connection acquires a WRITE table lock only that connection can read and write to the table.
--  Other connections cant read or write data until that lock is released by the connection.
LOCK TABLE messages WRITE;

UNLOCK TABLES;
