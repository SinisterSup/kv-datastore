# kv-datastore
## In-memory Key-Value DataStore

This is an in-memory Key-Value DataStore that performs operations on it based on certain commands and uses REST API for communication.

- Reads commands via HTTP REST API.
- Uses JSON encoding for requests and responses.
- Uses appropriate HTTP status codes for responses.

##### For detailed explanation with examples - refer [this](https://github.com/SinisterSup/kv-datastore/blob/main/Project%20details.pdf)
----

## Operations that are Implemented are:

### 1. SET :
   Writes the value to the datastore using the key and according to the specified parameters.    
   
   Pattern: `SET <key> <value> <expiry time>? <condition>?`   <br /> 
   
   `<key>`    <br /> 
   The key under which the given value will be stored.    
   `<value>`    <br /> 
   The value to be stored.     <br /> 
   `<expiry time>`   <br /> 
   Specifies the expiry time of the key in seconds.     <br /> 
   Must contain the prefix EX.     <br /> 
   This is an optional field,     <br /> 
   The field must be an integer value.         
   `<condition>`     <br /> 
   Specifies the decision to take if the key already exists.     <br /> 
   Accepts either NX or XX.     <br /> 
   NX -- Only set the key if it does not already exist.     <br /> 
   XX -- Only set the key if it already exists.     <br /> 
   This is an optional field. The default behavior will be to upsert the value of the key.    <br /> 
   
  #### - Use the Command of the form ->   
``` curl -X POST -H "Content-Type: application/json" -d '{"command": "SET hello world"}' http://localhost:8080 ``` 

  ---

### 2. GET :    
  Returns the value stored using the specified key.   
  Pattern: `GET <key>`    
  
  #### - Use the Command of the form ->   
``` curl -X GET -H "Content-Type: application/json" -d '{"command": "GET hello"}' http://localhost:8080 ``` 

  ---

### 3. QPUSH :    
  Creates a queue if not already created and appends values to it.    

  Pattern: `QPUSH <key> <value...>`  

  `<key>`    
  Name of the queue to write to.   
  `<value...>`      
  Variadic input that receives multiple values separated by space.    

  #### - Use the Command of the form ->   
```curl -X POST -H "Content-Type: application/json" -d '{"command":"QPUSH list_a a hola bella ciao"}' http://localhost:8080/``` 

  ---
     
### 4. QPOP :    
  Returns the last inserted value from the queue.    
         
  Pattern: `QPOP <key>`    
  
  `<Key>`    
  Name of the queue    
  
  #### - Use the Command of the form ->   
```curl -X GET -H "Content-Type: application/json" -d '{"command": "QPOP list_a"}' http://localhost:8080```

  ---

### 5. BQPOP :
  Blocking queue read operation that blocks the thread until a value is read from the queue.    
  The command must fail if multiple clients try to read from the queue at the same time.     

  Pattern: `BQPOP <key> <timeout>`   

  `<key>`
  Name of the queue to read from.    
  `<timeout>`   
  The duration in seconds to wait until a value is read from the queue.         
  The argument must be interpreted as a double value.       
  A value of 0 immediately returns a value from the queue without blocking.      
   
  #### - Use the Command of the form ->   
```curl -X GET -H "Content-Type: application/json" -d '{"command": "BQPOP list_a 0"}' http://localhost:8080```


----------------------------

### To Execute:- 
- Download or clone the repo    
- In the main directory (here named as kv-datastore) run the command --> ` go run main.go `    
(I'm assuming you have go installed on your machine and GOROOT, GOPATH variables are well managed and set.)
- Now open up a new Terminal to Test the APIs with the use of `"curl"` commands as I've suggested. 

#### Please feel free to raise an Issue if there is something I've missed out for. It would really help me out understand my mistakes!

# Thank You

