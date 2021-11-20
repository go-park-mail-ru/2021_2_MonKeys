#!/usr/bin/env tarantool

-- Настроить базу данных
box.cfg {
    listen = 3301
}


function init()
    s = box.schema.space.create('sessions', {if_not_exists=true})
    s:format({
        { name = 'cookie', type = 'string' },
        { name = 'user_id', type = 'unsigned' }
    })
    s:create_index('primary', { type = 'HASH', parts = {'cookie'}, if_not_exists = true })

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