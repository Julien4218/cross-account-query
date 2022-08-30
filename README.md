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
  nr_url: staging-api.newrelic.com/graphql
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
    nr_url: api.newrelic.com/graphql
    query: SELECT something else with NRQL
    select_fields:
      - is_using_something_else:IsUsingSomethingElse
```
