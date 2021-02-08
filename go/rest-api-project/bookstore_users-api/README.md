### Database setup steps

mysql commands to create the users table

```sh
CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY, first_name VARCHAR(45), last_name VARCHAR(45), email VARCHAR(45) NOT NULL, date_created DATETIME NOT NULL, status VARCHAR(45) NOT NULL, password VARCHAR(32) NOT NULL, UNIQUE INDEX email_unique (email ASC));
```
