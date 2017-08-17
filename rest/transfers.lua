-- 用合约来完成一个一对多转账系统
local L0 = require("L0")

-- 合约创建时会被调用一次，之后就不会被调用
function L0Init(args)
    L0.PutState("created", os.time())
    return true
end

-- 每次合约执行都调用
function L0Invoke(func, args)

    if("transfer" == func) then
        local args = {...}
        for m, n in ipairs(args) do
            local receiver = m
            local amount = tonumber(n)
            transfer(receiver, amount)
        end
    end

    return true
end

-- 查询
function L0Query(args)
    return "L0query ok"
end

function transfer(receiver, amount)
    L0.Transfer(receiver, amount)
end