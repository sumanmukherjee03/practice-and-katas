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
