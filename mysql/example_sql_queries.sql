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
insert into course_ids(select distinct cid from enrolled);





--  ##############################################################################################
--  ########################################## SUBQUERIES ########################################
--  ##############################################################################################

--  Subqueries are queries within queries which can be used within where or from clauses
select emp_name, dept, salary
from employees
where salary > (select avg(salary) from employees);

select emp_name, dept, salary
from employees
where salary > (select salary from employees where emp_name like 'John%');

select product_code, product_name, msrp from products
where product_code in (select product_code from orderdetails where price_each > 20 and price_each < 100);





--  ##############################################################################################
--  ############################################## UNION #########################################
--  ##############################################################################################

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
--  ############################################ RECURSION #######################################
--  ##############################################################################################
--  Recursive query
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


--  The first_value function gets the first entry for emp_name when the employees are ordered by salary.
--  So, all rows will show the same entry for the column highest_salary.
select emp_name, age, salary, first_value(emp_name) over (order by salary desc) as highest_salary from employees;


--  Show the name of the employee with highest salary in the department in which the current employee of the row belongs to.
select emp_name, age, salary, dept, first_value(emp_name) over (partition by dept order by salary desc) as highest_salary from employees;
