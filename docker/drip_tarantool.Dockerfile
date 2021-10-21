FROM tarantool/tarantool:latest

COPY ./docker/tarantool_scripts/dripapp.lua /opt/tarantool

EXPOSE 3301

CMD ["tarantool", "/opt/tarantool/dripapp.lua"]