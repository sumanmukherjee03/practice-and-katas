--  ##############################################################################################
--  ########################################## CASCADE DELETE ####################################
--  ##############################################################################################

CREATE TABLE buildings(
  building_no INT PRIMARY KEY AUTO_INCREMENT,
  building_name VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL
)

--  The ON DELETE CASCADE clause at the end is what causes the rooms to get deleted when the building is deleted

CREATE TABLE rooms(
  room_no INT PRIMARY KEY AUTO_INCREMENT,
  room_name VARCHAR(255) NOT NULL,
  building_no INT NOT NULL,
  FOREIGN KEY (building_no) REFERENCES buildings (building_no) ON DELETE CASCADE
)


--  As an admin this is a tip to find out which table is affected by ON DELETE CASCADE referential action
USE information_schema;
SELECT table_name FROM referential_constraints
  WHERE constraint_schema = 'buildings_db'
  AND referenced_table_name = 'buildings'
  AND delete_rule = 'CASCADE'
