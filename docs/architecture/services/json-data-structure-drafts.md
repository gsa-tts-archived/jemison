# **JSON data structures**

Related to: [https://github.com/orgs/GSA-TTS/projects/60/views/7?pane=issue\&itemId=97252640\&issue=GSA-TTS%7Cjemison%7C123](https://github.com/orgs/GSA-TTS/projects/60/views/7?pane=issue&itemId=97252640&issue=GSA-TTS%7Cjemison%7C123)

Below are the draft JSON data structures for our 3 sets of data:

* Search data  
* Search service data  
* Indexing information data

Specific questions we generated are in [Search 2.0 Data Planning](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=2103939643#gid=2103939643), in sheets titled "Indexing\_log\_questions" and "search\_guiding\_questions".

## **Reference**

Metric questions and planning are found at  [Search 2.0 Data Planning](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=0#gid=0).

## **Search data sample**

```json
{
        "query_uuid": "4f7ebebe-f2e4-11ef-bc21-df6c2092c112"
                // query ID is how we connect the search to the click metadata later
        "search_query": "passport",
                // resultsapi
        "domain": "www.usa.gov",
                // resultsapi
        "search": true,
        "referrer_url": "https://www.usa.gov",
                // result_data
        "query_duration": 321,
                // resultsapi
        "date_of_search": "2025-02-19",
                // resultsapi
        "list_of_results": [ { … }, { … } … ]
}

// This comes from the result_data component

{
        "query_uuid": "4f7ebebe-f2e4-11ef-bc21-df6c2092c112"
                // query ID is how we connect  click  to search metadata
        "click": true,
                // true if they click through (on a SERP page)
        "date_of_click": "2025-02-19",
        "click_url": "https://www.usa.gov/passport",
                // result_data
        "click_url_hostname": "www.usa.gov",
                // result_data
        "click_page_position": 1,
                // result_data
        "device_category": "Desktop",
                // result_data
        "search_type": "regular",
                // result_data
        "site_id": 9114
                // resultsapi
                // LVM: Consider renaming to be more descriptive
                // joined later with display_name, site_handle, agency

}
```

`result_data` is a new service. It provides the connecting glue between `resultsapi` (which will store a UUID for each query result) and the links that get clicked (which will need to ping an API endpoint with the UUID and any metadata we want to log about the query).

Provide on JS function that takes our results, and will attach handlers for each link. It can then do the work for our partners of registering a callback on each result link that will squirt a POST to our `results_data` API endpoint with the appropriate device, click, etc. metadata.

```json
{
  result_id: "4f7ebebe-f2e4-11ef-bc21-df6c2092c11d",
  ...
  results: [
  { 
    url: ["nasa.gov/spaceships](http://nasa.gov/spaceships")
    ...
  }
  ...
  ]
}
```

Original:

* From sheet: "search\_data\_fields" in [Search 2.0 Data Planning](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=0#gid=0).  
* Table format: [search\_data\_fields](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=620604994#gid=620604994)

## **Search service data sample**

Revised with summarized comments:
```json
[ // A variation of this, with additional fields developers may find helpful.
    {
        "date": "2025-02-19",
        "uptime": 0.99,
        "ram_usage": 0.101,
        "db_usage": 0.101
    }
]
```
Original:

* From sheet: "search\_service\_fields" in  [Search 2.0 Data Planning](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=1669337697#gid=1669337697)  
* Table format: [search\_service\_fields](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=1413418105#gid=1413418105)

## **Indexing information data sample**

Revised with summarized comments:

```json
{
        "domain_indexed": "www.usa.gov",
                // In the guestbook
        "url": "https://www.usa.gov/passport",
                // In the guestbook
        "canonical_url": "https://www.usa.gov/passport",
                // Could be logged, may not always be accurate. Is user input
        "status_code": 200,
        "status": "OK",
        "index_status": "indexed",
                // Do we need to extend the guestbook table with booleans for whether the `extract` service worked?
        "index_status_reason": null, 
        "index_status_date": "2025-02-19 15:08:40",
                // The status might get updated in case of a failure, and therefore would be more recent than the last_indexed date, which would only be updated when the indexing is successful.
        "last_successful_crawl": "2025-02-19 15:08:25",
                // This is set by `fetch` in the guestbook. Added _successful_ for clarity.
        "last_successful_index": "2025-02-19 15:08:40",
                // `extract` should be recorded in guestbook when successful. Added _successful_ for clarity.
        "last_modified": "2234-23-23 123:123:12",
                // From webserver. It is in the guestbook.
        "sha1": "3a423b32e2342...",
                // A hash of the fetched content. Possibly include metadata,  such as last_modified date from the webserver?
        "last_change_detected": "2025-02-19 1:36:34",
                // updated when hash changes
        "redirect_url": null,
                // MJ: "Need to figure out what `fetch` is doing now, and whether this makes sense to record. Because... we could get bounced through multiple redirects. What do we record?"
        "content_type": "html",
        "crawl_depth": 3,
                // MJ: This will be in the per-domain config. Add to the Guestbook? But may lead to redundant information in the system
        "first_contentful_paint": 0.44,
                // Javascript-rendered pages — future
        "response_time": 0.175,
                // Info obtained when fetched; potentially log to guestbook
        "payload_size": 9.7,
                // Currently logged, potentially under a different name. (ie. content length)
        "autothrottle_enabled": false
                // MJ: "In robots.txt. May want a whole different table, as it represents dynamic information we get from the user's site, and store to control our crawler."
    }

```

Original:

* From sheet: "indexing\_info\_data\_fields" in  [Search 2.0 Data Planning](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=1669337697#gid=1669337697)  
* Table format: [indexing\_info\_data\_fields](https://docs.google.com/spreadsheets/d/14KpCAXf-bVePMV-t56YTyPDvy3LuBf36KSw3QGZpgK4/edit?gid=2082126004#gid=2082126004)
