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

--  To substitute null values in the result set, the coalesce function can be used
select customerName, city, coalesce(state, 'N/A'), country from customers;
select id, release_date, coalesce(description, left(title, 150)), updated_at from movies;


--  Greatest and Least funcs usage along with ifnull
select company_id, least(ifnull(q1, 0), ifnull(q2, 0), ifnull(q3, 0), ifnull(q4, 0)) low, greatest(ifnull(q1, 0), ifnull(q2, 0), ifnull(q3, 0), ifnull(q4, 0)) high from revenues;

select * from orders where isnull(delivery_date);


--  Filter results based on aggregate computation
--  Here for example, you are aggregating the GPA but you cant use to filter results in your where clause
--  because it is part of the output.
--  So, to filter based on outputs of a selection use the having clause, because you can reference the output columns
select avg(s.gpa) as avg_gpa, e.cid
  from enrolled as e, student as s
  where s.sid = e.sid
  group by c.id
  having avg_gpa > 3.5

--  Get how much a product has sold
select productCode, sum(priceEach * quantityOrdered) as total
  from orderdetails
  group by productCode
  order by total desc;
