--  ##############################################################################################
--  ############################################# EVENTS #########################################
--  ##############################################################################################

SET GLOBAL event_scheduler = ON;
--  You can see the event_scheduler thread in the process list
SHOW PROCESSLIST;
SHOW EVENTS FROM testdb;


--  This creates a one time event
CREATE EVENT IF NOT EXISTS test_event_01
ON SCHEDULE AT CURRENT_TIMESTAMP
DO
  INSERT INTO messages(message, created_at) VALUES('Test event 01', NOW());

--  This creates a one time event that executes a minute after creation and is not dropped after completion
CREATE EVENT IF NOT EXISTS test_event_02
ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 MINUTE
ON COMPLETION PRESERVE
DO
  INSERT INTO messages(message, created_at) VALUES('Test event 02', NOW());

--  This creates a recurring event with a start and a stop time
CREATE EVENT IF NOT EXISTS test_event_03
ON SCHEDULE EVERY 1 MINUTE
STARTS CURRENT_TIMESTAMP + INTERVAL 1 HOUR
ENDS CURRENT_TIMESTAMP + INTERVAL 3 HOURS
DO
  INSERT INTO messages(message, created_at) VALUES('Test event - recurring', NOW());



SET GLOBAL event_scheduler = OFF;
