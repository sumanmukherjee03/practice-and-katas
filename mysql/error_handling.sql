--  ##############################################################################################
--  ######################################## ERROR HANDLING ######################################
--  ##############################################################################################

--  The syntax for exception handlers in mysql is
--  DECLARE action HANDLER FOR condition_value statement;
--  action can be either CONTINUE or EXIT
--  condition_value can be a mysql error code
--  or a standard SQLSTATE value or SQLWARNING, NOTFOUND, SQLEXCEPTION condition
--  or a named condition associated with a mysql error code or SQLSTATE value

--  For multiple error handlers, if there are multiple handlers that match the same error,
--  the precedence is mysql error code first, then SQLSTATE, then SQLEXCEPTION

declare continue handler for SQLEXCEPTION set has_error = 1;

delimiter $$
create procedure insert_supplier_product(in in_supplier_id int, in in_product_id int)
begin
  declare exit handler for 1062
  begin
    select concat('duplicate key (', in_supplier_id, ',', in_product_id, ') found') as message;
  end;

  insert into supplier_products (supplier_id, product_id) values (in_supplier_id, in_prod);
  select count(*) from supplier_products where supplier_id = in_supplier_id;
end$$
delimiter ;


--  Example of a named condition and error handling
drop procedure if exists test_proc;
delimiter $$
create procedure test_proc()
begin
  declare table_not_found_error condition for 1146;
  declare exit handler for table_not_found_error
    select 'please create table first' as message;
  select * from whatever;
end$$
delimiter ;



--  ##############################################################################################
--  ######################################## ERROR HANDLING ######################################
--  ##############################################################################################

--  sql allows you to raise errors with SIGNAL and RESIGNAL.
--  The syntax for SIGNAL is
--  SIGNAL SQLSTATE | condition_name;
--  SET condition_information_item_name_1 = value_1,
--      condition_information_item_name_1 = value_2, etc;

--  Here's an example of raising an unhandled generic user-defined exception via SIGNAL
delimiter $$
create procedure add_order_item(in orderno INT, in productcode VARCHAR(60), in qty int, in price double, in line_no int)
begin
  declare c int;
  select count(ordernumber) into c from orders where ordernumber = orderno;
  --  If the number of orders for that order number is not exactly 1 then this procedure raises an user defined sql exception.
  --  Error with SQLSTATE '45000' is a generic error for unhandled user-defined exception
  if (c != 1) then
    signal sqlstate '45000' set MESSAGE_TEXT = 'order not found in orders table';
  end if;
  --  Rest of the proceduce continues below .......
end$$
delimiter ;

--  Here's an example of resignal to provide a better error message
delimiter $$
create procedure divide(in num int, in deno int, out result double)
begin
  declare division_by_zero condition for sqlstate '22012';
  declare continue handler for division_by_zero resignal set MESSAGE_TEXT = 'Division by zero - denominator can not be zero';
  if deno = 0 then
    signal division_by_zero;
  else
    set result = num/deno;
  end if;
end$$
delimiter ;

call divide(11, 0, @result);
