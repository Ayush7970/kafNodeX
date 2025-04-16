-- benchmark/send_data.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["X-API-Key"] = "my_secret_api_key_12345"  -- Replace with your actual API key

-- Set a valid stream_id from a pre-created stream
local stream_id = "your_precreated_stream_id"  -- Replace with an actual stream_id from /stream/start
wrk.path = "/stream/" .. stream_id .. "/send"
wrk.body = '{"data": "test_data"}'
