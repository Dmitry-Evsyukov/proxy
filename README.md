# http-proxy

## Description:
http-proxy deployed on port 8080, proxying http and https requests
<br/><br/>
web api listening port 8000 and implementing sql-injection vulnerability scanner


### Examples proxy requests
* HTTP `curl -x http://localhost:8080 http://mail.ru/`
* HTTPS `curl -k https://mail.ru/ -x http://127.0.0.1:8080/ -vvv`

### Run

`cd deployment && docker compose up`

### API Description
1. `GET /api/v1/requests` – List of requests;
2. `GET /api/v1/requests/{id}` – Output 1 request;
3. `GET /api/v1/repeat/{id}` – Resubmit request;
4. `GET /api/v1/scan/{id}` – Request vulnerability scanner (sql-injection);

### SQL-Injection
SQL-Injection - add to the request headers/cookies/post_params/get_params single or double quote. If destination has vulnerability there is diff between old response and a new one
