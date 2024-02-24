create table pupils(
    id serial primary key,
    name varchar(255) not null,
    grade_id bigint default null references grade
);

create table grade(
    id serial primary key,
    mark int 
)