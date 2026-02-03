## 请求数据注意项

- overview的监控数据 5秒请求一次

## API 返回数据

### 服务器列表
- /api/v1/server

```
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "hostname": "hank",
      "ip_address": "127.0.0.1",
      "ssh_username": "root",
      "ssh_port": 22,
      "auth_type": "password",
      "status": ""
    }
  ]
}
```

- /api/v1/server/{id}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "hostname": "hank",
    "ip_address": "127.0.0.1",
    "ssh_username": "root",
    "ssh_port": 22,
    "auth_type": "password",
    "status": ""
  }
}
```

- /api/v1/monitor/stats/{server_id}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "cpu": {
      "cores": 4,
      "frequency": 2995.2,
      "model": "12th Gen Intel(R) Core(TM) i5-12490F",
      "perCoreUsage": [
        0.9999999980209395,
        0,
        0,
        0.9900990108139619
      ],
      "usage": 0.4975247522087254
    },
    "disk": {
      "available": 87585943552,
      "partitions": [
        {
          "available": 87253012480,
          "device": "/dev/sda3",
          "fsType": "ext4",
          "mountPoint": "/",
          "total": 104568143872,
          "usage": 12.051341493347392,
          "used": 11956019200
        },
        {
          "available": 332931072,
          "device": "/dev/sda2",
          "fsType": "ext4",
          "mountPoint": "/boot",
          "total": 473923584,
          "usage": 23.853331334782283,
          "used": 104292352
        }
      ],
      "total": 105042067456,
      "usage": 11.481411061384355,
      "used": 12060311552
    },
    "hostname": "hank",
    "loadAverage": {
      "load1": 0.66,
      "load15": 0.26,
      "load5": 0.56
    },
    "memory": {
      "available": 6462390272,
      "swapTotal": 4294963200,
      "swapUsed": 798720,
      "total": 8277491712,
      "usage": 21.928157745765194,
      "used": 1815101440
    },
    "timestamp": "2026-02-01T20:35:23.639426256+08:00",
    "topCPU": [
      {
        "cpuPercent": 0.31474959572051175,
        "createTime": 1769934144360,
        "memoryMB": 473.7734375,
        "memoryPercent": 0,
        "name": "node",
        "pid": 74343,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.20347408604448985,
        "createTime": 1769776993150,
        "memoryMB": 317.53515625,
        "memoryPercent": 0,
        "name": "Lingma",
        "pid": 12213,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.17257124467405072,
        "createTime": 1769949197200,
        "memoryMB": 27.078125,
        "memoryPercent": 0,
        "name": "agent",
        "pid": 77399,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.16459810607048897,
        "createTime": 1769949197100,
        "memoryMB": 34.60546875,
        "memoryPercent": 0,
        "name": "go",
        "pid": 77375,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.15252018709983428,
        "createTime": 1769949200110,
        "memoryMB": 39,
        "memoryPercent": 0,
        "name": "go",
        "pid": 77428,
        "status": "sleep"
      }
    ],
    "topMemory": [
      {
        "cpuPercent": 0.3147487290450982,
        "createTime": 1769934144360,
        "memoryMB": 473.7734375,
        "memoryPercent": 0,
        "name": "node",
        "pid": 74343,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.20347403609836995,
        "createTime": 1769776993150,
        "memoryMB": 317.53515625,
        "memoryPercent": 0,
        "name": "Lingma",
        "pid": 12213,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.04466560748435782,
        "createTime": 1769860710680,
        "memoryMB": 192.96875,
        "memoryPercent": 0,
        "name": "node",
        "pid": 54830,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.01882260492338619,
        "createTime": 1769863364250,
        "memoryMB": 89.58984375,
        "memoryPercent": 0,
        "name": "dockerd",
        "pid": 60397,
        "status": "sleep"
      },
      {
        "cpuPercent": 0.08751479146329225,
        "createTime": 1769860711130,
        "memoryMB": 85.07421875,
        "memoryPercent": 0,
        "name": "node",
        "pid": 54893,
        "status": "sleep"
      }
    ]
  }
}
```

- /api/v1/monitor/stats/io/{server_id}/{disk_name}

```

{
  "code": 0,
  "message": "success",
  "data": {
    "device": "sda3",
    "readBytes": 1172878336,
    "readCount": 36585,
    "readTime": 7440,
    "writeBytes": 19209805824,
    "writeCount": 304203,
    "writeTime": 159795
  }
}
```

- /api/v1/monitor/stats/io/{server_id}/all

```
{
  "code": 0,
  "message": "success",
  "data": {
    "device": "all",
    "readBytes": 2360252416,
    "readCount": 73827,
    "readTime": 14982,
    "writeBytes": 38420815872,
    "writeCount": 608550,
    "writeTime": 319649
  }
}
```

- /api/v1/monitor/stats/net/{server_id}/{interface_name}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "bytesRecv": 1536205557,
    "bytesSent": 411005376,
    "dropin": 0,
    "dropout": 0,
    "errin": 0,
    "errout": 0,
    "name": "ens33",
    "packetsRecv": 1671513,
    "packetsSent": 933090
  }
}
```

- /api/v1/monitor/stats/net/{server_id}/all

```
{
  "code": 0,
  "message": "success",
  "data": {
    "bytesRecv": 1957279815,
    "bytesSent": 838332602,
    "dropin": 0,
    "dropout": 0,
    "errin": 0,
    "errout": 0,
    "name": "all",
    "packetsRecv": 2963182,
    "packetsSent": 2225593
  }
}
```

- /api/v1/monitor/base/{server_id}/{page}/{count}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "collect_time": "2026-02-01T20:35:01.002387248+08:00",
        "cpu_usage": 0.25000000023283064,
        "disk_total": 105042067456,
        "disk_usage": 11.481368168093038,
        "disk_used": 12060266496,
        "id": 1,
        "memory_total": 8277491712,
        "memory_usage": 21.997038902018534,
        "memory_used": 1820803072
      }
    ],
    "page": 1,
    "size": 5,
    "total": 1
  }
}
```

- /api/v1/monitor/disk/{server_id}/{page}/{count}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "collect_time": "2026-02-01T20:35:01.004287753+08:00",
        "disk_name": "sda",
        "id": 1,
        "io_time": 0,
        "iops_in_progress": 0,
        "read_bytes": 1181325312,
        "read_count": 36965,
        "read_time": 7494,
        "weighted_io_time": 0,
        "write_bytes": 19208261632,
        "write_count": 304028,
        "write_time": 159766
      },
      {
        "collect_time": "2026-02-01T20:35:01.004287753+08:00",
        "disk_name": "sda1",
        "id": 2,
        "io_time": 0,
        "iops_in_progress": 0,
        "read_bytes": 425984,
        "read_count": 104,
        "read_time": 18,
        "weighted_io_time": 0,
        "write_bytes": 0,
        "write_count": 0,
        "write_time": 0
      },
      {
        "collect_time": "2026-02-01T20:35:01.004287753+08:00",
        "disk_name": "sda2",
        "id": 3,
        "io_time": 0,
        "iops_in_progress": 0,
        "read_bytes": 5608448,
        "read_count": 162,
        "read_time": 29,
        "weighted_io_time": 0,
        "write_bytes": 192512,
        "write_count": 26,
        "write_time": 18
      },
      {
        "collect_time": "2026-02-01T20:35:01.004287753+08:00",
        "disk_name": "sda3",
        "id": 4,
        "io_time": 0,
        "iops_in_progress": 0,
        "read_bytes": 1172878336,
        "read_count": 36585,
        "read_time": 7440,
        "weighted_io_time": 0,
        "write_bytes": 19208069120,
        "write_count": 304002,
        "write_time": 159747
      }
    ],
    "page": 1,
    "size": 5,
    "total": 4
  }
}
```

- /api/v1/monitor/net/{server_id}/{page}/{count}

```
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "bytes_recv": 1536650710,
        "bytes_sent": 411527031,
        "collect_time": "2026-02-01T20:40:01.011626245+08:00",
        "drop_in": 0,
        "drop_out": 0,
        "err_in": 0,
        "err_out": 0,
        "fifo_in": 0,
        "fifo_out": 0,
        "id": 2,
        "interface_name": "ens33",
        "packets_recv": 1673250,
        "packets_sent": 934686
      },
      {
        "bytes_recv": 1535913225,
        "bytes_sent": 410638186,
        "collect_time": "2026-02-01T20:35:01.004287753+08:00",
        "drop_in": 0,
        "drop_out": 0,
        "err_in": 0,
        "err_out": 0,
        "fifo_in": 0,
        "fifo_out": 0,
        "id": 1,
        "interface_name": "ens33",
        "packets_recv": 1670061,
        "packets_sent": 931741
      }
    ],
    "page": 1,
    "size": 5,
    "total": 2
  }
}
```