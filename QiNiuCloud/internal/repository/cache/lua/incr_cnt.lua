local key = KEYS[1]
--区分喜欢、下载和下载后关闭次数
local cntKey = ARGV[1]
local delta = tonumber(ARGV[2])
local exists=redis.call("EXISTS",key)
if exists==1 then
    redis.call("HINCRBY",key,cntKey,delta)
    return 1
else
    return 0
end