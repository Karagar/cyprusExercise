FROM mcr.microsoft.com/mssql/server:2019-latest

USER root

COPY . /

RUN chmod +x entrypoint.sh

CMD /bin/bash ./entrypoint.sh