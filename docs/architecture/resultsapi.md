# ResultsApi
## Purpose
The current search.gov system provides customers with a Search.gov Results API[^1]. In order to offer a smooth and easy transition Jemison will provide a a similar Results API so customers experience no change or system down time when switching over. 

### Results API
There are a number of required and optional parameters that can be found in other documentation[^1]. These parameters must be included for backwards compatibility. 

1. The API consumes the parameters to pass the query to the serve services. 
2. Once the API receives the data it packages and formats it into JSON to send back to the client. 

## ToDos
1. Testing and validation


## References
[^1]: https://open.gsa.gov/api/searchgov-results/
