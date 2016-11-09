FROM python:3.5.2-alpine

RUN adduser -Du 343 runner

USER runner
