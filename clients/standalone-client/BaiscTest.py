__author__ = 'jaunty'

import time
import random

import ccpmetrics

db_host = "10.86.1.170"
db_pass = "root"
db = "metrics"


def increment_test():
    lc = ccpmetrics.CCPMetrics(db_host, "root", db_pass, db)

    for i in range(20):
        lc.increment("counter2", random.randint(20,30), tags=["cool", "stuff"])
        lc.increment("counter2", random.randint(10,50))
        time.sleep(5)

    lc.stop()


def set_test():
    lc = ccpmetrics.CCPMetrics(db_host, "root", db_pass, db)

    for i in range(10):
        lc.set("set1", random.randint(1, 100))
        time.sleep(1)

    for i in range(10):
        lc.set("set2", str(random.randint(1, 100)))
        time.sleep(1)

    lc.stop()


def histogram_test():
    lc = ccpmetrics.CCPMetrics(db_host, "root", db_pass, db)

    for i in range(10):
        for j in range(100):
            lc.histogram("histogram1", random.randint(1, random.randint(1, 100)), tags=["cool", "histogram"])
            time.sleep(5)

    lc.stop()


def combined_test():
    lc = ccpmetrics.CCPMetrics(db_host, "root", db_pass, db)

    while True:
        if lc.get("counter").metric > 1000:
            lc.reset("counter")

        if lc.get("counter2").metric > 600:
            lc.reset("counter2")

        lc.increment("counter", random.randint(1, 10), tags=["users", "stuff"])
        lc.increment("counter2", random.randint(4, 12), tags=["users", "stuff"])
        lc.histogram("histogram1", random.randint(0, 100), tags=["requests", "histogram"])
        lc.gauge("gauge1", random.random(), tags=["load", "gauge"])

        time.sleep(1)


def combined_test2():
    lc = ccpmetrics.CCPMetrics(db_host, "root", db_pass, db)

    while True:
        flip = random.randint(1,10)
        if flip % 2 == 0:
            lc.increment("counter3", random.randint(1, 5), tags=["users", "stuff"])
        else:
            lc.increment("counter3", random.randint(4, 20), tags=["users", "stuff"])

        if lc.get("counter3").metric > 100:
            lc.reset("counter3")

        lc.gauge("gauge2", random.random()*random.randint(1,2), tags=["load", "gauge"])

        time.sleep(random.randint(1,3))

set_test()
# histogram_test()
#increment_test()
#combined_test2()
