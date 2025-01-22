# adding a new service

Adding a new service should generally follow steps that look like the following.

## the architecture

Our services communicate via queues. We use the [river](https://riverqueue.com) library to handle our queueing/messaging. There are three databases: one used by the queues (automatic; no work required on the developer's part), a "work" database (consider it a living scratchpad for services), and a "search" database, where our data ends up and is optimized for various use cases.

When there is data to store, we either use S3 or Postgres. 

## a new service, step-by-step

### create a new `cmd`

In the `cmd` directory, create a new folder for the service. We'll call this new example service `searchapi`. Because these will become golang paths, lets avoid using any symbols, etc. in the service names.

### create `main.go`

In the service folder, create a `main.go`.

```go
package main
	
import (
    "sync"
)

func main() {
  var wg sync.WaitGroup
  wg.Add(1)
  wg.Wait()
}
```

### add a `Makefile`

You can copy an existing file, or create a new one. At the least, it should provide `run`, `build`, and `clean`.

```
run: clean
  go run *.go

build: clean generate
	go build -buildvcs=false -o service.exe

clean:
	rm -f service.exe
```

All services are output as `service.exe` for consistency.

### add to the container stack

In the top-level `compose.yaml`, add a new entry for the service. You might begin by copying an existing service. Note that it will (possibly) need a port number (if it has any internal connections). Otherwise, it should be straight-forward:

```
searchapi:
  <<: *services-common
  image: jemison/dev
  # Simulate CF
  # https://stackoverflow.com/questions/42345235/how-to-specify-memory-cpu-limit-in-docker-compose-version-3
  deploy:
    resources:
      limits:
        memory: 64m
  build: 
    context: .
    dockerfile: ./cmd/searchapi/Dockerfile
  entrypoint: /home/vcap/app/cmd/searchapi/run.sh
  volumes:
    - type: bind
      source: .
      target: /home/vcap/app
  ports:
    - 10007:8888
  # https://docs.docker.com/compose/how-tos/startup-order/
  depends_on:
    nginx:
      condition: service_started
    jemison-queues-db:
      condition: service_healthy
    jemison-work-db:
      condition: service_healthy
  healthcheck:
    test: curl --fail http://searchapi:8888/heartbeat || exit 1
    interval: 60s
    timeout: 180s
    retries: 3
    start_period: 60s
  environment:
    ENV: "DOCKER"
    PORT: 8888
    DEBUG_LEVEL: debug
    GIN_MODE: debug
    SCHEDULE: ${SCHEDULE:-""}
```

Things to check:

* `depends_on`: If there are services that we need to wait for, put those here.
* `healthcheck`: this should be the same across the board
* `environment`: thsese want to be double-checked (this is a "to do"), but setting the `ENV` to `DOCKER` is a must. 

(The containerization/runtime configuration is still a bit in flux/wants discussion/design work. But, copy-pasta-ing this is probably fine for the moment.)

### test

At this point, test that the build works, and the service comes up. The example provided has an infinite wait. It probably will *not* respond to healthchecks... we'll add that in a moment.

### add configuration

Every service has configuration that must be defined.

In `config/services`, add a new service Jsonnet file.

E.g. in `searchapi.jsonnet`:

```jsonnet
local B = import 'base.libsonnet';
local service = 'searchapi';

local credentials = [
  [
    'port',
    { cf: 8080, container: 8888, localhost: 8888 },
  ],
];

local parameters = [
  [
    'debug_level',
    { cf: 'warn', container: 'debug', localhost: 'debug'},
  ],
] + B.parameters;

{
  creds:: [[service] + x for x in credentials],
  params:: [[service] + x for x in parameters],
  cf: B.params('credentials', 'cf', service, self.creds) +
      B.params('parameters', 'cf', service, self.params),
  container: { name: service } +
             B.params('credentials', 'container', service, self.creds) +
             B.params('parameters', 'container', service, self.params),
}
```

This is a good, basic config. It imports some common config, sets the debug level (which every service is expected to have), and that's it. This will automatically be slurped up by the build.

### add a healthcheck

We need to now add some common infrastructure to the service.

```
var ThisServiceName = "searchapi"

func main() {
  env.InitGlobalEnv(ThisServiceName)
  engine := common.InitializeAPI()

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
```

We may want to consider having some global constants for service names, instead of strings. However, there's no way to make that constant work across the config files and the applications, so... :shrug:. We use strings for the moment.

It should now be possible to stand up the stack with the new service, and have it respond to healthchecks. The common API initialization establishes a basic healthcheck for every single service. We need this so that Cloud.gov will be able to tell if our services are alive/responding.

### add to the deployment

You should, when you're ready, add to the TF deployment. We'll document this later (when we have a standardized TF deploy).

### add the queues

If the service communicates via queues, you need to add that in. There are two ways it might communicate:

1. As a worker
2. As a job creator

We'll handle each in turn.

#### add a job creator

If the service creates jobs for other services, there is common code for this.

In `main`:

```
var ThisServiceName = "searchapi"
var ChQSHP = make(chan queueing.QSHP)

func main() {
  env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

  engine := common.InitializeAPI()

	go queueing.Enqueue(ChQSHP)

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
```

The two lines do two things:

1. We create a channel global to the service called ChQSHP. This is short for "Channel for Queue name, Scheme, Host, and Port." A channel is like a wire; it lets you send data in one end, and at the other, a process will pick up the data and do stuff.
2. We pass one end of the channel to `Enqueue` in the `queueing` library. (This is an internal library that is part of Jemison.) 

Now, we can, anywhere in the service, send a message down this channel. This is how we enqueue new jobs for other services (or, oddly, even the service we're writing.)

```
		ChQSHP <- queueing.QSHP{
			Queue:      "entree",
			Scheme:     "https",
			Host:       "blogs.nasa.gov",
			Path:       "/astronauts",
		}
```

This says "create a QSHP data structure, and pack it with values. (The channel takes one value; to pass multiple values, we pass them in a struct.) Then, once packed, we send the value over the channel with the `<-` operator.

* https://go.dev/tour/concurrency/2
* https://gobyexample.com/channels

The work of sending the message to the queue (a DB table) is handled by the `Enqueue` function. If we add services, we need to extend the `Enqueue` function to handle them. It isn't magic, and it has some checking to make sure the string we pass (the service name) is one we know about. In other words, it serves as a gatekeeper to make sure we don't send messages off into the void.

#### as a job worker

The above is how we send jobs *to* the queue. What if we want to work jobs from the queue?

Add a function call to `main`:

```
var ThisServiceName = "searchapi"

func main() {
  env.InitGlobalEnv(ThisServiceName)
	InitializeQueues()

  engine := common.InitializeAPI()

	zap.L().Info("listening to the music of the spheres",
		zap.String("port", env.Env.Port))
	// Local and Cloud should both get this from the environment.
	http.ListenAndServe(":"+env.Env.Port, engine)
}
```

It needs to come after the env is initialized (or we don't have logging), and should probably come before most anything else (because we want to fail fast if we can't establish communication with the DB).

Now, in a file called `queues.go`, add that function. It probably looks like this to start:

```go
// The work client, doing the work of `SearchApi`
var SearchApiPool *pgxpool.Pool
var SearchApiClient *river.Client[pgx.Tx]

type SearchApiWorker struct {
	river.WorkerDefaults[common.SearchApiArgs]
}

func InitializeQueues() {
	queueing.InitializeRiverQueues()

	ctx, fP, workers := common.CommonQueueInit()
	SearchApiPool = fP

	// Essentially adds a worker "type" to the work engine.
	river.AddWorker(workers, &SearchApiWorker{})

	// Grab the number of workers from the config.
	SearchApiService, err := env.Env.GetUserService(ThisServiceName)
	if err != nil {
		zap.L().Error("could not SearchApi service config")
		log.Println(err)
		os.Exit(1)
	}

	// Work client
	SearchApiClient, err = river.NewClient(riverpgxv5.New(SearchApiPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			ThisServiceName: {MaxWorkers: int(SearchApiService.GetParamInt64("workers"))},
		},
		Workers: workers,
	})

	if err != nil {
		zap.L().Error("could not establish worker pool")
		log.Println(err)
		os.Exit(1)
	}

	// Start the work clients
	if err := SearchApiClient.Start(ctx); err != nil {
		zap.L().Error("workers are not the means of production. exiting.")
		os.Exit(42)
	}
}
```

Someday, it would be nice to make this boilerplate go away. However, the generics involved make standardizing it difficult. I haven't figured it out yet...

This code...

1. Initializes the queues. This happens in every service. 
2. We use a common init to get a DB context, a pool, and a list of workers. (The list is empty.)
3. We set a global (to the service) to the pool. We use it elsewhere.
4. Register the worker for this service with the library.
5. Get the configuration for this service.
6. Use the configuration value to create a work client, and assign the number of workers defined in the config file.
7. Start the workers.

We have to add the job arguments structure to our common infrastructure:

```go
type SearchApiArgs struct {
	Scheme    string `json:"scheme"`
	Host      string `json:"host"`
	Path      string `json:"path"`
}

func (SearchApiArgs) Kind() string {
	return "searchapi"
}
```

This is in `internal/common/types.go`. This datastructure is what gets turned into JSON and stuck in the DB as part of the queueing process. Most services in Jemison only pass URLs, but if you need to pass additional data, it will go here. This is how you get data from one service to the next.

We have one more thing to do: we have to define the workers.

In a file called `work.go`:

```go
func (w *FetchWorker) Work(ctx context.Context, job *river.Job[common.FetchArgs]) error {

  // ... possibly turn the arguments into a URL? 
	u := url.URL{
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

  // ... do stuff ...

  // ... enqueue something to another service?
	ChQSHP <- queueing.QSHP{
		Queue:  "extract",
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

  // Return `nil` if everything went well.
  // Returning an error anywhere in the work function will cause it to be requeued.
	return nil
}
```


## draw the owl

For those who are not familiar with the "draw the owl" meme:

https://www.reddit.com/r/funny/comments/eccj2/how_to_draw_an_owl/?rdt=34092

From here, you just have to draw the owl. :D 

Humor aside, the question is "what does your service do?" If this is the search API, then you're going to be adding code to implement a Gin API engine. It will probably **not** take jobs from the queue, because it is intended to receieve search queries, hit the `search` database, and return JSON data.

However... some questions:

1. How are we going to log data about searches? Will it... talk directly to a database? If so, we have some boilerplate for talking to the databases.
   1. Or, will we enqueue our data, and have a service that grabs all performance/app data, and stores it somewhere in the work DB? This way, all our services might have some internal metrics that they can simply enqueue for the `metrics` service to consume and store.
2. Do we ever need to control the search engine? For example, do we ever want to turn it off? Configure it? If so, then perhaps our `admin` component will enqueue messages for this service, and those will be "control signals." If that's the case, then we will need a work component, and it will watch for those messages.
   1. Why would we do this *this way*? Because we do *not* want to introduce the idea that services talk to each-other via API. We have queues. Use them.

There might be others. But, the point being: this is where we start drawing the owl.