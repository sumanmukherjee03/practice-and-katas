--  ##############################################################################################
--  #################################### FILTERING AND ORDERING ##################################
--  ##############################################################################################

--  This is going to produce different results each time because the order is not guaranteed
select sid, name from students where login like '%@cs.org' limit 20 offset 20
--  This is going to produce consistent results because the order is guaranteed unless ofcourse new rows get added
select sid, name from students where login like '%@cs.org' order by sid asc, name asc limit 20 offset 20
