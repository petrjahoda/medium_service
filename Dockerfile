FROM scratch
COPY /css /css
COPY /html html
COPY /js js
COPY /linux /
CMD ["/medium_service"]