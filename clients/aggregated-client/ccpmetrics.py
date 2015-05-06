__author__ = 'jaunty'

import httplib
import time
import json
from socket import gethostname
from time import gmtime, strftime

class CCPMetrics(object):
    def __init__(self, host, port, service=None):
        self.host = host
        self.port = port
        self.service = service

    def gauge(self, metric, value, tags=None, sample_rate=1):
        self._write(metric, value, "gauge", tags, sample_rate)

    def increment(self, metric, value=1, tags=None, sample_rate=1):
        self._write(metric, value, "counter", tags, sample_rate)

    def decrement(self, metric, value=1, tags=None, sample_rate=1):
        self._write(metric, -value, "counter", tags, sample_rate)

    def histogram(self, metric, value, tags=None, sample_rate=1):
        self._write(metric, value, "histogram", tags, sample_rate)

    # def timing(self, metric, value, tags=None, sample_rate=1):
    #     pass

    def set(self, metric, value, tags=None, sample_rate=1):
        self._write(metric, value, "set", tags, sample_rate)

    def _write(self, metric, value, metric_type, tags=[], sample_rate=1):
        http_serv = httplib.HTTPConnection(self.host, self.port)
        hostname = gethostname()

        output = json.dumps({
            "metric": metric,
            "host": hostname,
            "timestamp": strftime("%a, %d %b %Y %H:%M:%S +0000", gmtime()),
            "type": metric_type,
            "value": value,
            "sampling": sample_rate,
            "tags": tags,
            "service": self.service
            })

        http_serv.request('POST', '/metrics', output)
        print output

import random

for i in range(10000):
    tags = dict()
    tags["a"] = "1"
    tags["b"] = "2"
    metrics = CCPMetrics("localhost", 8082)
    metrics.gauge("foo1", random.randint(0, i/5), tags)
    metrics.histogram("bar1", random.randint(0, i/5), tags)

    if i % 3 == 0:
            time.sleep(0.05)