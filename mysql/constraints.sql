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
