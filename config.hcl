prometheus_url = "http://localhost:9090"

metrics {
  go_gc_duration_seconds {
    query = "go_gc_duration_seconds"
    asg = "testing"
    namespace = "tray/testing"
    unit = "None"
    interval = 1
  }

  sum_go_gc_duration_seconds {
    query = "sum(go_gc_duration_seconds)"
    asg = "testing"
    namespace = "tray/testing"
    unit = "None"
    interval = 1
  }
}
