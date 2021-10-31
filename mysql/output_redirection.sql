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
