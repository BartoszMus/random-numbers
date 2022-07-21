# random-numbers
A simple REST service in Go which performs concurrent requests to random.org API. The service calculates standard deviation of the drawn integers, standard deviation of sum of all sets and returns them in JSON format.

## Start
It can be launched by running the following command in the project directory
```
$ go run .
```
 
 It can also be started via docker after building image
 ```
  docker build --tag backend .
  docker run -p 8000:8000 backend
  
 ```
  ##Usage
  After it is launched, the data can be accessed in a browser by going to
  ```
  http://localhost:8000/random/mean?requests=4&length=3
  ```
  Number of requests and their lengths can be changed by editing according parameters
