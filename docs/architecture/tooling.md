# tooling

* **We deploy to Cloud.gov**. This is a Cloud Foundry-based environment. 
* **Never touch production**. We never log into `production` to manually do *anything*, unless it is an emergency/we are debugging something that cannot be debugged any other way.
* **Nothing is manual**. We trigger `production` deploys manually. That's it. 
* **We work in Go**. Our primary language for Jemison is Go. If we must introduce another language, our primary preference will be Python. We avoid `bash` at all costs.
* **We are open**. Our work is in the public domain. We leverage the best open tools we can. We do not reach for proprietary solutions.
* **We automate via GH Actions**. To interact with the CF/Cloud.gov environment, GH Actions are our preferred source of automations.
* **We use S3 and Postgres**. We think about our problem before we reach for new tools. Every new piece of infrastructure has a cost, both in dollars and maintenance/complexity.

