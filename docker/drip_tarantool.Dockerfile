FROM tarantool/tarantool:latest

COPY ./docker/tarantool_scripts/dripapp.lua /opt/tarantool

CMD ["tarantool", "/opt/tarantool/dripapp.lua"]