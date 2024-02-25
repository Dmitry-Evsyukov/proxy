create table request (
    id  serial  primary key,
    data json not null
);

create table response (
    id  serial  primary key,
    data json not null
);