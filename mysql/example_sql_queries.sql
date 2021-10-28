--  ##############################################################################################
--  ########################################## CONSTRAINTS #######################################
--  ##############################################################################################

--  Example of table with various constraints and an index
create table employees(
  id int not null,
  email varchar(255) not null,
  first_name varchar(255),
  last_name varchar(255),
  age int,
  dept varchar(255) default 'operations',
  primary key (id),
  unique(email),
  check(age > 20)
);
create index employees_name_index on employees(first_name, last_name);
desc employees;
insert into employees(id, email, first_name, last_name, age) values(1, 'john.doe@cs.org', 'John', 'Doe', 23);

show indexes from employees;
drop index employees_name_index on employees;

--  Example of primary key auto increment contraint example
create table posts(
  id int auto_increment primary key,
  content text,
  category_id int
) engine = INNODB;


create table races (
  id int not null auto_increment primary key,
  race_name varchar(30) not null
) engine = INNODB;

--  The on update cascade clause updates the race_id in characters table if id of a record is changed in the races table that is referenced in the characters table
--  The on delete restrict clause prevents deletion of a record from races table if it is referenced as a foreign key in any record of characters table
create table characters (
  id int not null auto_increment primary key,
  character_name varchar(50) not null,
  race_id int not null,
  index idx_race (race_id),
  constraint fk_character_race foreign key (race_id) references races(id) on update cascade on delete restrict
) engine = INNODB;


--  ##############################################################################################
--  ########################################## AGGREGATE FN ######################################
--  ##############################################################################################
--  Aggregation : AVG(col), MIN(col), MAX(col), COUNT(col), SUM(col)

select count(*) as cnt from student where login like '%@cs.org';
select count(distinct login) from student where login like '%@cs.org';

--  Get average GPA and number of students whose are in CS department
select avg(gpa), count(sid) from student where login like '%@cs.org';

--  Ex : Get the average GPA of students enrolled in each course
--  We need to make use of a group by here to aggregate on subsets of students enrolled in a course.
--  Also, in the group by clause we need to group by course id because non-aggregated values in select output clause
--  must appear in group by clause
select avg(s.gpa), e.cid from enrolled as e, student as s where e.sid = s.sid group by e.cid;
select avg(s.gpa), e.cid, e.cname from enrolled as e, student as s where e.sid = s.sid group by e.cid, e.cname;

--  Filter results based on aggregate computation
--  Here for example, you are aggregating the GPA but you cant use to filter results in your where clause
--  because it is part of the output.
--  So, to filter based on outputs of a selection use the having clause, because you can reference the output columns
select avg(s.gpa) as avg_gpa, e.cid
  from enrolled as e, student as s
  where s.sid = e.sid
  group by c.id
  having avg_gpa > 3.5




--  ##############################################################################################
--  ######################################### STRING MATCHING ####################################
--  ##############################################################################################
--  Strings use single quotes
--  like clause uses '%' matches any substring like a wildcard, even empty string
--  like clause uses '_' matches any one character
--  You can use string functions like substring, upper, lower etc in output section like select or even in predicates like the where clause
select * from enrolled as e where e.cid like '15-%';
select substring(name, 0, 5) as abbrev_name from student where sid = 55348;
select * from student as s where upper(s.name) like 'KAN%';

--  String concatenation can be done with concat function or `||` operator
select name from student where login = concat(lower(name), '@cs.org');


-- Get the current date and time
select now();
select current_timestamp;

--  Get the day field from a date, for example for 2021-10-21 get 21
select extract(day from date('2021-10-21'));

--  Get the number of days since the begining of the year
--  This one works in postgres
select date('2021-10-21') - date('2021-01-01') as days;
--  This one works in mysql
select datediff(date('2021-10-21'),date('2021-01-01')) as days;



--  ##############################################################################################
--  ###################################### OUTPUT REDIRECTION ####################################
--  ##############################################################################################

--  Output redirection
--  Send query results into other queries or tables

--  Here the database already knows the schema of the course_ids table
select distinct cid into course_ids from enrolled;

--  Here the database creates a new table for course_ids and it infers the schema based on the output of the select query
create table course_ids(select distinct cid from enrolled);

--  In this case you insert into an existing table
--  However, it is important to keep in mind that output redirection as inserts into another table may fail if the table has constraints
insert into course_ids(select distinct cid from enrolled);





--  ##############################################################################################
--  #################################### FILTERING AND ORDERING ##################################
--  ##############################################################################################

--  This is going to produce different results each time because the order is not guaranteed
select sid, name from students where login like '%@cs.org' limit 20 offset 20
--  This is going to produce consistent results because the order is guaranteed unless ofcourse new rows get added
select sid, name from students where login like '%@cs.org' order by sid asc, name asc limit 20 offset 20





--  ##############################################################################################
--  ########################################## SUBQUERIES ########################################
--  ##############################################################################################

--  Subqueries are queries within queries which can be used within where or from clauses
--  The inner query can reference outputs from the outer query

select emp_name, dept, salary
from employees
where salary > (select avg(salary) from employees);

select emp_name, dept, salary
from employees
where salary > (select salary from employees where emp_name like 'John%');

select product_code, product_name, msrp from products
where product_code in (select product_code from orderdetails where price_each > 20 and price_each < 100);


--  Find all the sudent names that are enrolled in courses with ids between 15 to 245
select name from students where sid in (select sid from enrolled where cid between 15 and 245)
--  Find all the courses that have no students enrolled in them
select * from courses where not exists (select * from enrolled where courses.cid = enrolled.cid)





--  ##############################################################################################
--  ############################################## JOINS #########################################
--  ##############################################################################################

--  Inner join is for finding intersection of 2 sets

--  Find players that are registered both in cricket and football matches
select c.match_id, c.player, f.match_id, f.player from cricket as c inner join football as f on c.player = f.player

--  Here's another example of an inner join but using the `using` keyword
select productcode, productname, textdescription from products inner join productlines using(productline);

--  This is an example of 3 way inner join
select o.ordernumber, o.status, p.productname, sum(od.quantityordered * od.priceofeachproduct) as revenue
from orders as o inner join orderdetails as od on o.ordernumber = od.ordernumber
inner join products as p on p.productcode = od.productcode
group by ordernumber;

--  A left join returns all the rows from the left table and the matching rows from the right table

--  In this example below there will be some rows where ordernumber and status will be empty
--  because those customers have no orders
select c.customernumber, c.customername, o.ordernumber, o.status from customers as c left join orders as o on c.customernumber = o.customernumber;

--  Using the information above, this query finds customers who have not placed any orders
select c.customernumber, c.customername, o.ordernumber, o.status from customers as c left join orders as o on c.customernumber = o.customernumber where ordernumber is null;


--  A right join returns all the rows from the right table and the matching rows from the left table

--  The example below has all values selected from employees table which is the right table.
--  There will be rows where customername and customer phone will be null.
--  This means not all employees have sold stuff to customers as sales reps
select c.customername, c.phone, e.employeenumber, e.email
from customers as c right join employees as e on e.employeenumber = c.salesrepnumber
order by employeenumber;


--  A self join is a join of a table with itself

--  This result set returns a manager with multiple employees reporting to him/her.
select concat(m.lastname, ', ', m.firstname) as manager, concat(e.lastname, ', ', e.firstname) as employee
from employees as e inner join employees as m
on m.employee_number = e.reports_to
order by manager;


--  Full outer join
--  If a flavour of sql does not support full outer join, the best way to achieve it is
--  by combining the result sets of a left outer join and right outer join via the union operator
select c.customername, o.ordernumber from customers as c
left join orders as o on c.customernumber = o.customernumber
union
select c.customername, o.ordernumber from customers as c
right join orders as o on c.customernumber = o.customernumber


--  When 2 tables have the same structure for the columns you select, and you want to join the rows, ie
--  you want rows from both the tables show up in 1 result, you use UNION.
--  The keyword UNION eliminates duplicate rows, but UNION ALL keeps duplicate rows in the result set.
--  You got to select the same number of columns from both the tables
select SalesOrderId, OrderDate, SubTotal from sales.StoreSales
  union select SalesOrderId, OrderDate, SubTotal from sales.EcommSales;

--  If the columns dont match you can add a placeholder blank/empty column
select SalesOrderId, OrderDate, SubTotal, SalesPersonId from sales.StoreSales
  union select SalesOrderId, OrderDate, SubTotal, '' from sales.EcommSales;






--  ##############################################################################################
--  ###################################### COMMON TABLE EXPRESSIONS ##############################
--  ##############################################################################################

--  CTE is an alternative to temp tables and views. It provides a temp table to work with for the lifetime of a query
--  You can also bind output columns to names for use in the select clause
with cteName (col1, col2) as
(
  select 1, 2
) select col1 + col2 from cteName

--  Find student with the highest id that is enrolled in at least one course
with cteSource (maxId) as (
  select max(sid) from enrolled
) select name from students, cteSource where students.sid = cteSource.maxId


-- Simple example of a recursive query with CTE that produces an output of numbers from 1 to 10
--  The recursive keyword here may be required in postgres and may not be required in mysql
with recursive cteSource (counter) as
(
  (select 1)
  union all
  (select counter + 1 from cteSource where counter < 10)
) select * from cteSource;

--  A more complex recursive query
--  ParentEmployeeKey here represents the employee key of the manager
--  If the ParentEmployeeKey of someone is NULL, it means he/she is the CEO or Boss
with EmployeeHierarchy as
(
  --  Get the Boss/CEO
  --  We are union ing 2 separate result sets. So make sure they have the same number of cols
  --  This first query is the anchor query where the level field is getting materialized first
  select EmployeeKey, FirstName, LastName, ParentEmployeeKey, 1 as level
    from DimEmployee
    where ParentEmployeeKey is null
  union
  --  Get the rest of the employees
  --  Here we are using a self join essentially
  --  This is the part with the recursive query
  select e.EmployeeKey, e.FirstName, e.LastName, e.ParentEmployeeKey, eh.level + 1 as level
  from DimEmployee e
    inner join EmployeeHierarchy eh
    on eh.EmployeeKey = e.ParentEmployeeKey
    where e.ParentEmployeeKey is not null
) select * from EmployeeHierarchy order by level;




--  ##############################################################################################
--  ####################################### STORED PROCEDURES ####################################
--  ##############################################################################################

--  List stored procs
show procedure status like '%pattern%'

--  Drop stored procs
drop procedure if exists abc;

--  Example of a simple procedure
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




--  ##############################################################################################
--  ############################################# VIEWS ##########################################
--  ##############################################################################################

--  Example of creating a view
create view customer_details
as
select customerName, phone, city from customers;

--  Example of invoking a view
select * from customer_details;


create view prod_desc
as
select productName, quantity, msrp, text_description
from products as p inner join product_lines as pl
on p.product_line = pl.product_line;

select * from prod_desc;

--  To rename a view use this command
rename table prod_desc to product_description;

--  To list all views use the command below
show full tables where table_type = 'VIEW';

--  To delete a view
drop view product_description;




--  ##############################################################################################
--  ############################################ WINDOW FN #######################################
--  ##############################################################################################

--  This is an example of a report of each farmer's orange production alongside the total orange production of the year
--  Here, the over clause constructs a static window that includes all the records returned by the query, ie, all the records for year 2020
select farmer_name, kilos_produced, sum(kilos_produced) over() total_produced from orange_production where crop_year = 2020;

--  Sliding or dynamic window frames mean the window of records can be different for each row returned by a query.
--  Moreover, the window is created based on the current row in the query, so the rows in the window can change when the current row changes.
--  The over(partition by orange_variety, crop_year) clause in the example below creates a dynamic sliding window by grouping all records with the same value in the orange_variety and crop_year columns
--  So, the columns total_same_variety_year in the result will return sum of kilos produced for valencia oranges in 2020, lets say.
--  Jake | Valencia | 2019 | 107 | 12000
--  Jake | Valencia | 2020 | 110 | 10000
--  John | Valencia | 2019 | 90  | 12000
--  John | Valencia | 2020 | 100 | 10000
--  Jill | Golden   | 2020 | 75  | 9800
--  Jill | Golden   | 2020 | 85  | 10500
select farmer, orange_variety, crop_year, kilos_produced, sum(kilos_produced) over(partition by orange_variety, crop_year) as total_same_variety_year from orange_production;

--  Example : problem is to find the combined salary of each department
--  So, we plan to partition our table by department and get the total salary
select emp_name, age, department, sum(salary) over (partition by dept) as total_salary_by_dept from employees;

select row_number() over (order by salary) as row_num, emp_name, salary from employees order by salary;

--  row_number can also be used to track duplicates
--  Here row_number is acting like a count.
--  It counts the number of records with the same id and name.
--  So, if there are 3 students with id 7 and name 'John', this col will display 3
select st_id, st_name, row_number() over (partition by st_id, st_name order by st_id) as row_num from students;


--  This example will rank the rows based on the value of field1
--  So, if multiple rows have the same value for field1, they will have the same rank
select field1, rank() over (order by field1) as test_rank from demo;

--  Find the student with the highest grade for each course
--  The inner query here groups tuples by cid and then sorts a group by grade
--  The ranking table is a temp table that exists only for the lifetime of this query
select * from (
  select *, rank() over (partition by cid order by grade asc) as rank from enrolled
) as ranking
where ranking.rank = 1;


--  The first_value function gets the first entry for emp_name when the employees are ordered by salary.
--  So, all rows will show the same entry for the column highest_salary.
select emp_name, age, salary, first_value(emp_name) over (order by salary desc) as highest_salary from employees;


--  Show the name of the employee with highest salary in the department in which the current employee of the row belongs to.
select emp_name, age, salary, dept, first_value(emp_name) over (partition by dept order by salary desc) as highest_salary from employees;
