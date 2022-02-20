FROM ilyagunagdimaev/dependencies:latest AS builder

WORKDIR /dripapp

CMD /dripapp/auth_service 