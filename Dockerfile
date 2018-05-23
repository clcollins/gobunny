FROM scratch
LABEL maintainer="Chris Collins <collins.christopher@gmail.com>"

COPY pkg/* /
# COPY client.crt /
# COPY client.key /
CMD [ "/gobunny" ]
