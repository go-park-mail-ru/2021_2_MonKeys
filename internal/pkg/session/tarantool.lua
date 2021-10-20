#!/usr/bin/env tarantool

-- Настроить базу данных
box.cfg {
    listen = 3301
}

-- При поднятии БД создаем спейсы и индексы
box.once('init', function()
    if not box.space.sessions then 
        box.schema.space.create('sessions')
        box.space.sessions:create_index('primary',
        { type = 'HASH', parts = {1, 'unsigned'}})
        box.space.sessions:create_index('secondary',
        { type = 'HASH', parts = {2, 'string'}})

        print('create sessions space')
    end
end)