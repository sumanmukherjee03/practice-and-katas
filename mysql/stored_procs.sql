--  ##############################################################################################
--  ####################################### STORED PROCEDURES ####################################
--  ##############################################################################################

--  List stored procs
show procedure status like '%pattern%'

--  Drop stored procs
drop procedure if exists abc;

--  Example of a simple procedure
drop procedure if exists top_players;
delimiter &&
create procedure top_players()
begin
  select name, country, goals from players where goals > 5;
end &&
delimiter ;

call top_players();


--  Example of procedure that takes 1 input
delimiter //
create procedure SortBySalary(IN num_recs int)
begin
  select name, age, salary
  from emp_details
  order by salary desc limit num_recs;
end //
delimiter ;

call SortBySalary(10);


--  Example of procedure that takes 2 inputs
delimiter $$
create procedure update_salary(IN temp_name varchar(20), IN new_salary float)
begin
  update emp_details set salary = new_salary where name = temp_name;
end $$
delimiter ;

call update_salary('John', 80000);


--  Example of procedure that has an output
delimiter $$
create procedure count_employees(OUT Total_Emps int)
begin
  select count(name) into Total_Emps from emp_details where sex = 'F';
end $$
delimiter ;

call count_employees(@F_emps);
select @F_emps as female_emps;


--  You can also modify a variable passed to the stored proc using the INOUT parameter type
delimiter $$
create procedure SetCounter(INOUT counter int, IN inc int)
begin
  set counter = counter + inc;
end$$
delimiter ;

set @counter = 1; -- set the counter variable to 1
call SetCounter(@couter, 3);
select @counter; -- 4


--  You can use conditions inside stored procs.
--  The IF statement has three forms: simple IF-THEN statement, IF-THEN-ELSE statement, and IF-THEN-ELSEIF-ELSE statement.
delimiter $$
create procedure customer_level(IN customer_id int, OUT customer_level varchar(20))
begin
  declare credit decimal(10, 2) default 0;
  select credit_limit into credit from customers where id = customer_id;
  if credit > 50000 then
    set customer_level = 'platinum';
  elseif credit > 10000 and credit <= 50000 then
    set customer_level = 'gold';
  else
    set customer_level = 'silver';
  end if;
end$$
delimiter ;

call customer_level(11, @level);
select @level;


--  Stored procs support case statements
delimiter $$
create procedure get_shipping_time(IN customer_id int, shipping_info varchar(50))
begin
  declare customer_country varchar(100);
  select country into customer_country from customers where id = customer_id;
  case customer_country
  when 'US' then
    set shipping_info = '2-day shipping';
  when 'CA' then
    set shipping_info = '3-day shipping';
  else
    set shipping_info = '5-day shipping';
  end case;
end$$
delimiter ;

call get_shipping_time(9, @shipping);
select @shipping;


--  In case you dont want to throw any error if no cases match you can have an empty begin+end in the else clause
--  The case statements can also come with multiple conditions instead of just matching
delimiter $$
create procedure get_delivery_status(IN customer_id int, OUT delivery_status varchar(100))
begin
  declare waiting_days int default 0;
  select datediff(est_delivery_date, shipped_date) into waiting_days from orders where customer_id = customer_id;
  case
  when waiting_days = 0 then
    set delivery_status = 'ontime';
  when waiting_days >= 1 and waiting_days < 3 then
    set delivery_status = 'arriving';
  when waiting_days >= 3 and waiting_days < 5 then
    set delivery_status = 'late';
  when waiting_days >= 5 and waiting_days < 10 then
    set delivery_status = 'verylate';
  else
    set delivery_status = 'lost';
end$$
delimiter ;

call get_delivery_status(7, @delivery);
select @delivery;

--  This is an example of looping inside a procedure
delimiter $$;
create procedure example_loop;
begin
  declare x int;
  declare str varchar(100);
  set x = 1;
  set str = '';
  --  Here the loop is labeled as looper
  looper: LOOP
    if x > 10 then
      LEAVE looper;
    end if;
    set x = x + 1;
    if (x mod 2) then
      iterate looper;
    else
      set str = concat(str, x, ',');
    end if;
  end LOOP;
  select str;
end$$
delimiter ;

call example_loop();


--  This is an example of a while loop in stored proc
delimiter $$
create procedure load_calendars(IN start_date DATE, IN no_days INT)
begin
  declare counter INT default 1
  declare dt DATE default start_date
  while counter <= no_days do
    insert into calendars (fulldate, day, month, quarter, year) values (dt, extract(day from dt), extract(month from dt), extract(quarter from dt), extract(year from dt));
    set counter = counter + 1;
    set dt = date_add(dt, interval 1 day)
  end while;
end$$
delimiter ;

call load_calendars('2021-10-31', 30);



-- This is an example of a repeat/until loop
delimiter $$
create procedure repeat_demo()
begin
  declare counter INT default 1;
  declare result VARCHAR(100) default '';
  repeat
    set result = concat(result, counter, ',');
    set counter = counter + 1;
  until counter >= 10;
  end repeat;
  select result;
end$$
delimiter ;
call repeat_demo();



--  An important thing about LEAVE clause is that just like it can be used to exit labelled loops/while etc
--  it can also be used to leave a stored proc if it labelled
DELIMITER $$
CREATE PROCEDURE sp_name()
sp: BEGIN
  [label:] WHILE search_condition DO
    IF condition THEN
      --  This leaves the loop
      LEAVE [label];
    END IF;
  END WHILE [label];

  --  Other statements

  IF condition THEN
    --  This leaves the stored proc
    LEAVE sp;
  END IF;
END$$
DELIMITER ;



--  ##############################################################################################
--  ############################### STORED PROCEDURE SECURITY CONTEXT ############################
--  ##############################################################################################

--  In this sample stored proc below, a user dev@localhost if allowed access execute stored procs in the current db
--  can invoke the stored proc InsertMessage. And even if the user might not have access to the messages table
--  he/she can insert data into it via this stored proc because the execution of the stored proc will run in the
--  context of the root user. Ofcourse, it's creation has to be done by the root user.
DELIMITER $$
CREATE DEFINER = root@localhost PROCEDURE InsertMessage(msg VARCHAR(100))
SQL SECURITY DEFINER
BEGIN
  INSERT INTO messages(message) VALUES(msq);
END$$
DELIMITER ;

GRANT EXECUTE ON testdb.* TO dev@localhost;


--  However this stored proc will run with the context of the invoker even though the root user created and defined it.
--  This means the dev user above wont be able to execute this because he/she does not have access to the messages table.
DELIMITER $$
CREATE DEFINER=root@localhost PROCEDURE UpdateMessage(
    msgId INT,
    msg VARCHAR(100)
)
SQL SECURITY INVOKER
BEGIN
    UPDATE messages
    SET message = msg
    WHERE id = msgId;
END$$
DELIMITER ;
