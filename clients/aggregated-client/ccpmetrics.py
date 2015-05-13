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
        self._write_metric(metric, value, "gauge", tags, sample_rate)

    def increment(self, metric, value=1, tags=None, sample_rate=1):
        self._write_metric(metric, value, "counter", tags, sample_rate)

    def decrement(self, metric, value=1, tags=None, sample_rate=1):
        self._write_metric(metric, -value, "counter", tags, sample_rate)

    def histogram(self, metric, value, tags=None, sample_rate=1):
        self._write_metric(metric, value, "histogram", tags, sample_rate)

    # def timing(self, metric, value, tags=None, sample_rate=1):
    #     pass

    def set(self, metric, value, tags=None, sample_rate=1):
        self._write_metric(metric, value, "set", tags, sample_rate)

    def event(self, title, text, alert_type=None, aggregation_key=None,
              source_type_name=None, priority=None,
              tags=None, hostname=None, date_happened=None):

        http_serv = httplib.HTTPConnection(self.host, self.port)

        output = json.dumps({
            "name": title,
            "text": text,
            "host": hostname,
            "alerttype": alert_type,
            "priority": priority,
            "timestamp": strftime("%a, %d %b %Y %H:%M:%S +0000", gmtime()),
            "AggregationKey": aggregation_key,
            "SourceType": source_type_name,
            "tags": tags,
            })

        http_serv.request('POST', '/events', output)

    def _write_metric(self, metric, value, metric_type, tags=[], sample_rate=1):

        http_serv = httplib.HTTPConnection(self.host, self.port)
        hostname = gethostname()

        output = json.dumps({
            "name": metric,
            "host": hostname,
            "timestamp": strftime("%a, %d %b %Y %H:%M:%S +0000", gmtime()),
            "type": metric_type,
            "value": value,
            "sampling": sample_rate,
            "tags": tags,
            "service": self.service
            })

        http_serv.request('POST', '/metrics', output)
