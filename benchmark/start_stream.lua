-- benchmark/start_stream.lua
wrk.method = "POST"
wrk.path = "/stream/start"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["X-API-Key"] = "my_secret_api_key_12345"  -- Replace with your actual API key


