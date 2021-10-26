#!/usr/bin/env tarantool

-- Настроить базу данных
box.cfg {
    listen = 3301
}


function init()
    if box.space.sessions then box.space.sessions:drop() end

    s = box.schema.create_space('sessions')
    s:format({
        { name = 'cookie', type = 'string' },
        { name = 'user_id', type = 'unsigned' }
    })
    s:create_index('primary', { type = 'HASH', parts = {'cookie'} })

    print('create space sessions')
end

function new_session(cookie, user_id)
    print('received data', user_id)

    box.space.sessions:insert{cookie, user_id}

    print('insert cookie', cookie)

    return 'lol'
end

function check_session(cookie)
    local session = box.space.sessions:select{cookie}
    print('found session', user_id)
    return session
end

function delete_session(cookie)
    local  droped_session = box.space.sessions:delete{cookie}
    print('delete session', deleted_session)
    return droped_session
end