Create database temporary;
Go
USE temporary;
Go
Create table test (id int, title nvarchar(64));
Insert into test(id, title) values (1, 'One'), (2, 'Two');
