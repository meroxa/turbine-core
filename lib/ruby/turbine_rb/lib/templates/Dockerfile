FROM ruby:3.1.2
WORKDIR /app
COPY . .
RUN bundle install
USER nobody
ENTRYPOINT ["bundle", "exec", "turbine-function"]
