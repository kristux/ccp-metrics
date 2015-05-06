__author__ = 'jaunty'

from influxdb import InfluxDBClient
import redis


class RedisWriter(object):
    def __init__(self, host, username, password, database=0, port=6379):
        self.redis_client = redis.StrictRedis(host, port, database)

    def write(self, metric):
        json_body = [
            {
                "name": str(metric.name),
                "tags": metric.tags,
                #"time": int(metric.timestamp),
                "columns": metric.columns(),
                "points": [
                    metric.values()
                ]
            }
        ]

        print json_body
        #self.redis_client.rpush('metrics', json_body)


class InfluxWriter(object):

    def __init__(self, host, username, password, database, port=8086):
        self.client = InfluxDBClient(host, port, username, password, database)

    def write(self, metric):

        json_body = [
            {
                "name": str(metric.name),
                "tags": metric.tags,
                #"time": int(metric.timestamp),
                "columns": metric.columns(),
                "points": [
                    metric.values()
                ]
            }
        ]

        print json_body

        self.client.write_points(json_body)
