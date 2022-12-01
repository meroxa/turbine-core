# frozen_string_literal: true

require "rubygems"
require "bundler/setup"
require "turbine_rb"

class MyApp
  def call(app)
    database = app.resource(name: "demopg")

    # ELT pipeline example
    # records = database.records(collection: 'events')
    # database.write(records: records, collection: 'events_copy')

    # procedural API
    records = database.records(collection: "events")

    # This register the secret to be available in the turbine application
    app.register_secrets("MY_ENV_TEST")

    # you can also register several secrets at once
    # app.register_secrets(["MY_ENV_TEST", "MY_OTHER_ENV_TEST"])

    # Passthrough just has to match the signature
    processed_records = app.process(records: records, process: Passthrough.new)
    database.write(records: processed_records, collection: "events_copy")

    # out_records = processed_records.join(records, key: "user_id", window: 1.day) # stream joins

    # chaining API
    # database.records(collection: "events").
    #   process_with(process: Passthrough.new).
    #   write_to(resource: database, collection: "events_copy")
  end
end

# might be useful to signal that this is a special Turbine call
class Passthrough < TurbineRb::Process
  def call(records:)
    puts "got records: #{records}"
    # to get the value of unformatted records, use record .value getter method
    # records.map { |r| puts r.value }
    #
    # to transform unformatted records, use record .value setter method
    # records.map { |r| r.value = "newdata" }
    #
    # to get the value of json formatted records, use record .get method
    # records.map { |r| puts r.get("message") }
    #
    # to transform json formatted records, use record .set methods
    # records.map { |r| r.set('message', 'goodbye') }
    records
  end
end

TurbineRb.register(MyApp.new)
