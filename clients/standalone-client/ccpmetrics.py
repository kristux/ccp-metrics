__author__ = 'jaunty'

import threading
import time
import Queue
from metrics import Counter, Metric, Set, Histogram
from output import RedisWriter, InfluxWriter


class CCPMetrics(object):

    def __init__(self, host, username, password, database, port=8086, buffer_size=256,
                 write_frequency=5):
        self.writer = InfluxWriter(host, username, password, database, port=8086)
        #self.writer = RedisWriter(host, username, password)
        self.running = True
        self.live_metrics = dict()
        self.write_queue = Queue.Queue()
        self.write_frequency = write_frequency

        self.queue_thread = threading.Thread(target=self.queue_thread)
        self.queue_thread.start()

    def increment(self, metric_name, value=1, tags=None):
        try:
            self.live_metrics[metric_name].metric += value
        except (KeyError, AttributeError):
            self.live_metrics[metric_name] = Counter(metric_name, tags)
            self.live_metrics[metric_name].metric += value
        except TypeError:
            print metric_name + " is not a counter"

    def decrement(self, metric_name, value=1, tags=None):
        try:
            self.live_metrics[metric_name].metric -= value
        except (KeyError, AttributeError):
            self.live_metrics[metric_name] = Counter(metric_name, tags)
            self.live_metrics[metric_name].metric -= value

    def gauge(self, metric_name, value, tags=None):
        new_gauge = Metric(metric_name)
        new_gauge.metric = value
        self.write_queue.put(new_gauge, tags)

    def set(self, metric_name, value, tags=None):
        try:
            self.live_metrics[metric_name].metric.add(value)
        except (KeyError, AttributeError):
            self.live_metrics[metric_name] = Set(metric_name, tags)
            self.live_metrics[metric_name].metric.add(value)
        except AttributeError or TypeError:
            print metric_name + " is not a set"

    def histogram(self, metric_name, value, tags=None):
        try:
            self.live_metrics[metric_name].metric.append(value)
        except (KeyError, AttributeError):
            self.live_metrics[metric_name] = Histogram(metric_name, tags)
            self.live_metrics[metric_name].metric.append(value)
        except AttributeError or TypeError:
            print metric_name + " is not a histogram"

    def queue_thread(self):
        while self.running:
            time.sleep(self.write_frequency)

            for metric_key in self.live_metrics:
                self.writer.write(self.live_metrics[metric_key])

    def reset(self, metric_name):
        del self.live_metrics[metric_name]

    def stop(self):
        self.running = False

    def get(self, metric_name):
        return self.live_metrics[metric_name]