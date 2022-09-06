# cross-account-query

A command line tool that execute multiple NRQL query accross environment and accounts.

Base query fetches the initial dataset. Column queries adds scalar data for each row in the initial dataset.

Output is displayed as CSV

Config example:

```yaml
base:
  api_key: NRAK-ABC**********
  account_id: 123456789
  region: staging
  query: SELECT something with NRQL
  select_fields:
    - account_id:AccountID
    - customer_user_id:UserID
    - timestamp
    - timestamp:CopyDateUTC

columns:
  -
    api_key: NRAK-ABC**********
    account_id: 987654321
    region: us
    query: SELECT something else with NRQL
    select_fields:
      - is_using_something_else:IsUsingSomethingElse
```

## batch support

config column query can be executed in batch mode, which execute much faster instead of an N+1 execution plan
to leverage batching, you'll want to do the following:
- use the `can_batch: true` option
- write the query using a `WHERE` clause in the form of `...WHERE field1 IN (env::my_field)...` so that all the field values previously queried can be injected
- add the replaced field in the query select, so the results can be merged with the previous results

For example:
```yaml
base:
  api_key: NRAK-ABC**********
  account_id: 123456789
  region: staging
  query: SELECT something with NRQL
  select_fields:
    - account_id:AccountID
    - customer_user_id:UserID
    - timestamp
    - timestamp:CopyDateUTC

columns:
  -
    api_key: NRAK-ABC**********
    account_id: 987654321
    region: us
    can_batch: true
    query: FROM MyTable SELECT uniqueCount(1) as 'has_data', min(customer_user_id) as 'customer_user_id' WHERE customer_user_id in (env::customer_user_id) SINCE 1 week ago facet nrAccountId as 'account_id'
    select_fields:
      - is_using_something_else:IsUsingSomethingElse
```


# Run

```bash
go test -v ./ ./...
go run . config.yml
```
