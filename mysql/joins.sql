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
