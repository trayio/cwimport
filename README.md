### cwimport

Send data from Prometheus to CloudWatch for auto scaling.

### Warning

Data received from Prometheus seems to be model.Vector no matter what the query is.
However this is probably wrong assumption that's why we log if data is of different type.
