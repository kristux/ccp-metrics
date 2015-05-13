import random
import time
from ccpmetrics import CCPMetrics

metrics = CCPMetrics('localhost', 8082)

for i in range(100000):
	tags = dict()
	tags["host"] = "host1"
	tags["service"] = "httpd"
	metrics.gauge("foo3", random.randint(1, 100), tags)
	metrics.histogram("foo4", random.randint(1, 100), tags)
	metrics.increment("foo5", random.randint(1, 100), tags)
	metrics.event("event1", "a wild event occured", "severe", "event1key", "client", "high", tags=tags) 
	time.sleep(0.01)
