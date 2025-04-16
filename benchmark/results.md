# Benchmark Results

## `/stream/start` Endpont
- **Concurrency**: 1000 connections
- **Requests/sec**: 119382.89
- **Latency (avg)**: 10.08 ms
- **Transfer/sec**: 21.63 MB
- **Total Requests**: 3592741
- **Total Data Transferred**: 651.00 MB
- **Socket Errors**:
  - Connect: 0
  - Read: 3314
  - Write: 0
  - Timeout: 0
- **Non-2xx or 3xx responses**: 3592741

---

## `/stream/{stream_id}/send` Endpoint
- **Concurrency**: 1000 connections
- **Requests/sec**: 119921.22
- **Latency (avg)**: 11.49 ms
- **Transfer/sec**: 21.73 MB
- **Total Requests**: 3608935
- **Total Data Transferred**: 653.93 MB
- **Socket Errors**:
  - Connect: 0
  - Read: 3272
  - Write: 0
  - Timeout: 0
- **Non-2xx or 3xx responses**: 3608935
