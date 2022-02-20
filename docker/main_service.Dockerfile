FROM ilyagunagdimaev/dependencies:latest AS builder

WORKDIR /dripapp

CMD /dripapp/main_service