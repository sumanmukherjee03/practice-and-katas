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
