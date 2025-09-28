local key = KEYS[1]
local likeCntKey = ARGV[1]
local fieldDownloadCntKey = ARGV[2]
local closeAfterDownloadedCntKey = ARGV[3]
local likeCntValue = tonumber(ARGV[4])
local fieldDownloadCntValue = tonumber(ARGV[5])
local closeAfterDownloadedValue = tonumber(ARGV[6])
local exists=redis.call("EXISTS",key)
if exists==1 then
    return -1
else
    redis.call('HINCRBY', key,likeCntKey, likeCntValue)
    redis.call('HINCRBY', key,fieldDownloadCntKey, fieldDownloadCntValue)
    redis.call('HINCRBY', key,closeAfterDownloadedCntKey, closeAfterDownloadedValue)
    return 1
end