# databases

The engine is supported by multiple databases. We do this because different databases may want/need to scale differently in the future, and we are taking a preemptive step in that direction.

## unique data values

There are a small number of unique data representations in our application that merit note.

* the [domain64](domain64.md) representation for domain names. This is a 64-bit integer representation of domain names that we use for search optimization and table partitioning.

## queues

The queues database serves only one purpose: to handle the queues used by Jemison.

Our queueing system gets hit hard, and therefore we do all of that work on one database. Further, the table migrations are automanaged by the [river](https://riverqueue.com/) library. Keeping it separate both protects the queue tables as well as any operational tables we create in the application.

## work

The "work" database is where application tables specific to the processing of data live. 

Read more about the [tables and their roles in the work database](databases/work.md).

## search

The `search` database holds our data pipelines and the tables that get actively searched. 

Read more about the [tables and their roles in the search database](databases/search.md).
