#! /usr/bin/env tarantool

local box = require('box')
local msgpack = require('msgpack')

box.once("bootstrap", function()
    box.schema.space.create('kv_storage')

    box.kv_storage:format({
        { name = 'key', type = 'str' },
        { name = 'value', type = 'map' },
    })

    box.space.kv_storage:create_index('primary', {
        parts = { 'key' },
    })

    msgpack.cfg{
        encode_invalid_as_nil = true,
    }
end)
