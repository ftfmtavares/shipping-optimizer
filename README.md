# shipping-optimizer
An exercise implementation of shipping packages optimizer based on different packing sizes and ordered amount



## Makefile Commands
make run:          # Run Locally<br>
make test:         # Unit Tests<br>
make build:        # Production Build<br>
make vet:          # Code Quality Checks<br>
make docker-build: # Generate Docker Image (requires Docker)<br>
make docker-up:    # Run locally in a Docker Container (requires Docker)<br>
make docker-down:  # End Docker Container Execution (requires Docker)<br>



## API Endpoints Examples
Retrieve the package sizes of a given product<br>
GET 'http://localhost:8080/product/1/packsizes'<br>
Response:<br>
{<br>
    "pid": 1,<br>
    "packs": [ 23, 31, 53 ]<br>
}<br>


Setup the package sizes of a given product<br>
POST 'http://localhost:8080/product/1/packsizes'<br>
Payload:<br>
{<br>
    "packs": [23, 31, 53]<br>
}<br>
Response:<br>
{<br>
    "pid": 1,<br>
    "packs": [ 23, 31, 53 ]<br>
}<br>


Setup the package sizes of a given product<br>
GET 'http://localhost:8080/product/1/shipping-calculation?order=500000'<br>
Response:<br>
{<br>
    "order": 500000,<br>
    "packs": [<br>
        {<br>
            "packsize": 23,<br>
            "quantity": 2<br>
        },<br>
        {<br>
            "packsize": 31,<br>
            "quantity": 7<br>
        },<br>
        {<br>
            "packsize": 53,<br>
            "quantity": 9429<br>
        }<br>
    ],<br>
    "packscount": 9438,<br>
    "total": 500000,<br>
    "excess": 0<br>
}<br>



## UI Demonstration Page
http://localhost:8080/

![alt text](misc/image.png)



## Live Demo
The application is running online on Render:<br>
https://shipping-optimizer-latest.onrender.com/<br>