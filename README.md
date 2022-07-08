## Http Traffic Counter

Counts the requests made to the webserver and persists them as JSON file
and shows the client the number of requests made in the moving window (by default
it is 1 minute)

```go
const DefaultExportFileName = "requests.json"
const DefaultMovingWindow = 1 * time.Minute
```

Also, by default the data will be persisted to `requests.json` in project folder

### The used method

Aimed to have 100% accuracy in data. 
Different approach would bet to have different goroutine that works
in cron style and every second or so calculates the amount of requests and persists.
But this approach would not guarantee 100% accuracy as between the cron
executions, there may be more requests. But this approach would be much
faster than the current implementation