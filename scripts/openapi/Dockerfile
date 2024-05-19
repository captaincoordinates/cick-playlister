FROM python:3.12-alpine3.19

RUN pip install openapi-spec-validator

WORKDIR /openapi

COPY internal/docs/openapi.yml .
CMD ["openapi-spec-validator", "openapi.yml"]
