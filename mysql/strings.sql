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
