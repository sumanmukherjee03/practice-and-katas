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
