#! /usr/bin/env tarantool

-- Import required modules
local box = require('box')
local msgpack = require('msgpack')

-- Runs only once during the initialization.
box.once("bootstrap", function()
    --- Create a space.
    box.schema.space.create('kv_storage')

    --- Define schema.
    box.kv_storage:format({
        { name = 'key', type = 'str' },
        { name = 'value', type = 'map' },
    })

    --- Add primary index.
    box.space.kv_storage:create_index('primary', {
        parts = { 'key' },
    })

    --- MsgPack serialization option.
    msgpack.cfg{
        encode_invalid_as_nil = true,
    }
end)
