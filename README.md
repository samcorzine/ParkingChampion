# Parking Champion

Parking Champion is an api designed to allow the storage and query of parking rates. The set of current rates can be updated and queried via a RESTful endpoint, and users can find the applicable rate for given starting and end times via a getRate endpoint.

#### Running the API

1. Make sure docker is properly installed `https://docs.docker.com/install/`
2. Clone the repo
3. cd into cloned directory
4. `docker build -t parking . && docker run --rm -p 8080:8080 parking`

### Usage

API will be initialized with no rates, so first off you'll want to post a valid rates document to the /rates endpoint. Feel free to use the one stored at testing/testing.json. The following curl command will do the intitial post for you if run from the top-level directory of the cloned repo.

`curl -i -X POST localhost:8080/rates  -H "Content-Type: text/xml"   --data "@testfiles/testing.json"`

From there rates can be queried via the getRate endpoint.

`curl "localhost:8080/getRate?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00"`

Endpoints are fully documented in the swagger.yaml file.

### Metrics

Metrics are exposed at /metrics, courtesy of the SpotHero Tools Library for Go.
`https://github.com/spothero/tools`
