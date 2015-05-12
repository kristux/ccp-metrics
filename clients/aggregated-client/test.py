import random
import time
from ccpmetrics import CCPMetrics

metrics = CCPMetrics('localhost', 8082)

for i in range(100):
	tags = dict()
	tags["host"] = "host1"
	tags["service"] = "httpd"
	metrics.gauge("foo3", random.randint(1, 100), tags)
	metrics.histogram("foo4", random.randint(1, 100), tags)
	metrics.increment("foo5", random.randint(1, 100), tags)
	time.sleep(1)
