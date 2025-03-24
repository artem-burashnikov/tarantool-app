#! /usr/bin/env tarantool

local box = require('box')
local msgpack = require('msgpack')

box.once("bootstrap", function()
    box.schema.user.create('storage', { password = 'sesame' })

    box.schema.user.grant('storage', 'super', nil, nil)

    box.schema.space.create('default_space')

    box.space.default_space:create_index('primary', {
        type = 'TREE',
        parts = {1, 'unsigned'},
    })

    msgpack.cfg{
        encode_invalid_as_nil = true,
    }
end)
