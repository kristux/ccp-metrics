__author__ = 'jaunty'

import time
import numpy
import json


class Metric(object):
    metric = 0

    def __init__(self, metric_name, tags=None, timestamp=None):
        self.name = metric_name
        self.tags = tags

        if timestamp is None:
            self.timestamp = time.time()
        else:
            self.timestamp = timestamp

    def columns(self):
        return ["count"]

    def values(self):
        return self.metric

    def tags(self):
        return self.tags

    def values(self):
        return [self.metric]

    def columns(self):
        return ["metric"]


class Counter(Metric):
    def __init__(self, metric_name, tags=None, timestamp=None):
        super(Counter, self).__init__(metric_name, tags, timestamp)

    def columns(self):
        return ["count"]

    def increment(self, value=1):
        self.metric += value

    def decrement(self, value=1):
        self.metric -= value


class Set(Metric):
    metric = set()

    def __init__(self, metric_name, tags=None, timestamp=None):
        super(Set, self).__init__(metric_name, tags, timestamp)

    def values(self):
        return json.dumps(list(self.metric))

    def columns(self):
        return ["set"]


class Histogram(Metric):
    metric = list()

    def columns(self):
        return ["count", "min", "max", "mean", "std-dev",
                "50-percentile", "75-percentile",
                "95-percentile", "99-percentile"]

    def values(self):
        count = len(self.metric)
        minimum = min(self.metric)
        maximum = max(self.metric)
        mean = numpy.mean(self.metric)
        stddev = numpy.std(self.metric)
        percentile50 = numpy.percentile(self.metric, 50)
        percentile75 = numpy.percentile(self.metric, 75)
        percentile95 = numpy.percentile(self.metric, 95)
        percentile99 = numpy.percentile(self.metric, 99)

        return [count, minimum, maximum, mean, stddev, percentile50,
                percentile75, percentile95, percentile99]