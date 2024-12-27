# migrate

Migrate is not a service. It is a small go program that runs database migrations.

It:

1. reads the `dbmate` migrations for the `work` and `search` datbases, 
2. stores them in an embedded filesystem at compile time, and
3. runs those stored/compiled-in migrations using `dbmate` as a Go library at runtime.

In may systems, we would ship the migrations as text, and run them (perhaps by running `dbmate` via a `bash` script) on the host. Instead, we compile the migrations into the `migrate` app, and run that app at deploy time.

This reduces the number of moving pieces in the system, and keeps a critical management action in Go.