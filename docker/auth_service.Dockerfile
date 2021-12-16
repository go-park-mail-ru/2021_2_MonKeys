FROM dependencies AS builder

WORKDIR /dripapp

CMD /dripapp/auth_service 