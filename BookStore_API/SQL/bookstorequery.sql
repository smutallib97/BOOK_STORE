show databases;
create database bookstore;
use bookstore;
create table users (user_ID INT unsigned not null, userName varchar(30) not null, email varchar(30) not null, password varchar(30), primary key (user_ID));
select * from users;
insert into users (user_ID, userName, email, password) values (1, 'smutallib', 'smutallib97@gmail.com', 'pass@123'),
(2, 'dhiraj', 'dhiraj@gmail.com', 'pass12'), (3, 'raushan', 'thor@gmail.com', 'pass@456');
create table books (Book_ID INT unsigned not null, Book_Title varchar(30) not null, Book_Author varchar(30) not null, Book_Price INT unsigned not null, primary key (Book_ID));
select * from books;
insert into books (Book_ID, Book_Title, Book_Author, Book_Price) values (1, 'A Better India: A Better World', 'Narayana Murthy', 266), 
(2, 'An Introduction to Dreamland', 'Bhagat Singh', 49), (3, 'Harry Potter', 'J.K.Rowling', 255);
create table order_details (Order_ID INT unsigned not null, Book_Name varchar(30) not null, Address varchar(45) not null, Amount double not null, primary key (Order_ID));
select * from order_details;
insert into order_details (Order_ID, Book_Name, Address, Amount) values (1, 'Harry Potter', 'Buldhana', 255);
select * from order_details;

