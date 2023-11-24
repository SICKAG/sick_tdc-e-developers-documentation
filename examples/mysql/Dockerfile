FROM x11tete11x/arm32v7-mysql

COPY init.sh .

RUN sed -i '$ d' entrypoint.sh 
RUN echo './init.sh &' >> /entrypoint.sh
RUN echo 'exec "$@"' >> /entrypoint.sh