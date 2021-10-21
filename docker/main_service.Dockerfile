FROM dependencies AS builder

WORKDIR /app

COPY main_service /app/

CMD /app/main_service