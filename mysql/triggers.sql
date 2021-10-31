--  ##############################################################################################
--  ############################################ TRIGGERS ########################################
--  ##############################################################################################
create table students (st_roll int, age int, name varchar(30), marks float);

--  Example of creating a before insert trigger
delimiter $$
create trigger set_default_marks
before insert on students
for each row
if new.marks < 0 then set new.marks = 40;
end if; $$
delimiter ;

--  The insert statement will call the trigger
insert into students values(501, 10, 'Jane', 55.0),
  (502, 10, 'John', -1);

--  To delete a trigger
drop trigger set_default_marks;
