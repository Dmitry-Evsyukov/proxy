create table request (
    id  serial  primary key,
    scheme text,
    method text,
    url text,

    headers json,
    get_params json,
    cookies json,
    post_params json
);

create table response (
    id  serial  primary key,
    request_id int references request(id),
    code int,
    message text,
    headers json,
    body text
);
