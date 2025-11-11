# shipping-optimizer
An exercise implementation of shipping packages optimizer based on different packing sizes and ordered amount.  
<br>
The application is running online on Render:  
https://shipping-optimizer-latest.onrender.com/  
<br>

---

## Quick Start
  
#### Prerequisites
- Go 1.25+  
(or)  
- Docker
<br><br>

#### Clone Project
```sh
git clone -v https://github.com/ftfmtavares/shipping-optimizer
cd shipping-optimizer
```
<br>

#### Run
```sh
make run
```
<br>

#### Build
```sh
make build
# result: ./bin/shipping-optimizer
```
<br>
The aplication requires environment variables to be set before running the binary:  
- SERVER_ADDRESS (e.g. localhost)
- SERVER_PORT (e.g. 8080)
<br>

#### Run tests and coverage
```sh
make test
```
<br>

#### Build Docker Image:
```sh
make docker-build
```
<br>

#### Run Docker Container:
```sh
make docker-up
```
<br>

#### Stop Docker Image:
```sh
make docker-down
```
<br>

---

## API
A minimal demonstration UI is served at root:
http://localhost:8080/  
<br>
![alt text](misc/image.png)
<br><br>

#### Product Packages Size Configuration Set
- POST /product/{pid}/packsizes  
  Command:
```sh
curl -X POST -H "Content-Type: application/json" \
  -d '{"packs":[23,31,53]}' \
  http://localhost:8080/product/1/packsizes
```
  Response example:  
```json
{
    "pid": 1,
    "packs": [ 23, 31, 53 ]
}
```
<br>

#### Product Packages Size Configuration Read
- GET /product/{pid}/packsizes  
  Command:
```sh
curl -s http://localhost:8080/product/1/packsizes
```
  Response example:  
```json
{
    "pid": 1,
    "packs": [ 23, 31, 53 ]
}
```
<br>

#### Order Shipping Calculation
- GET /product/{pid}/shipping-calculation?order={qty}  
  Command:
```sh
curl -s "http://localhost:8080/product/1/shipping-calculation?order=500000"
```
  Response example:  
```json
{
    "order": 500000,
    "packs": [
        {
            "packsize": 23,
            "quantity": 2
        },
        {
            "packsize": 31,
            "quantity": 7
        },
        {
            "packsize": 53,
            "quantity": 9429
        }
    ],
    "packscount": 9438,
    "total": 500000,
    "excess": 0
}
```
<br>

#### Validation rules and limits
- pid = valid and non negative integer
- qty = valid and non negative integer (max 10B units)
- package size = non negative integer
<br><br>

---