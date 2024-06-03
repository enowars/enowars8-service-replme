import multiprocessing

worker_class = "uvicorn.workers.UvicornWorker"
workers = min(4, multiprocessing.cpu_count() * 2 + 1)
bind = "0.0.0.0:8000"
timeout = 90
keepalive = 3600

