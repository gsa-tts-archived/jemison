# tools

## wipe-queues.bash

Run this while the stack is running to clobber all of the `river-*` tables.

This is useful for clearing the state of a crawl, and making sure you have a clean slate. Could even be used in production in an emergency.

It does wipe out all of the queues, however. 

## reset-guestbook.bash

Does a drop on `guestbook` and `hosts`, and truncates the migrations.

A full reset on the work DB, basically.

## check-for-redirects.py

Takes a list of domains, like

```
www.sigpr.gov
erdc.hpc.mil
www.americorps.gov
www.cldp.doc.gov
frtib.gov
bep.gov
...
```

and checks to see what files have redirects. Writes the file `redirects.csv` to the local folder. That file has the form

```
url,resolved,status
www.sigpr.gov,https://www.sigpr.gov/,200
erdc.hpc.mil,https://erdc.hpc.mil/,301
erdc.hpc.mil,https://centers.hpc.mil/users/erdcdsrc.html,200
www.americorps.gov,https://www.americorps.gov/,200
www.cldp.doc.gov,https://www.cldp.doc.gov/,301
www.cldp.doc.gov,https://cldp.doc.gov/,200
frtib.gov,https://frtib.gov/,301
frtib.gov,https://www.frtib.gov/,200
bep.gov,https://bep.gov/,302
bep.gov,https://www.bep.gov/,200
...
```

Note that some URLs show up more than once. This is because the script is printing the redirect history. In this example, `www.cldp.doc.gov` ultimately redirects to `cldp.doc.gov`, which means we want to use the domain that has a `200` status in our scheduling/crawling configurations.

## build-schedule.ipynb

Takes a file of the form

```
espanol.mentalhealth.gov,46
cetsound.noaa.gov,47
www.research.iowa-city.med.va.gov,47
aas.gsa.gov,47
```

and the redirects file, merges them, and then generates a schedule file for consumption by `entree`.

```
{
  "weekly": [
    {
      "scheme": "https",
      "host": "www.9-11commission.gov",
      "path": "/"
    },
    {
      "scheme": "https",
      "host": "access.afsc.noaa.gov",
      "path": "/"
    },

    ...
```

We can adjust how we generate this schedule; for now, it is using numeric thresholds to try and give us a reasonable number of pages to crawl each (week/month/quarter).