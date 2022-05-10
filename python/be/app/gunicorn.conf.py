# -*- coding: utf-8 -*-

import multiprocessing
import os
from distutils.util import strtobool

bind = os.getenv('BE_ADDR', '0.0.0.0:8002')
accesslog = '-'
access_log_format = "%(h)s %(l)s %(u)s %(t)s '%(r)s' %(s)s %(b)s '%(f)s' '%(a)s' in %(D)sµs"
#workers = int(os.getenv('BE_WORKERS', multiprocessing.cpu_count() * 2))
workers = int(os.getenv('BE_WORKERS', 2))
threads = int(os.getenv('BE_THREADS', 1))
reload = bool(strtobool(os.getenv('BE_RELOAD', 'false')))