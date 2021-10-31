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

