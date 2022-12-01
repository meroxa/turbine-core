require 'json'
require 'hash_dot'

module TurbineRb
  class Record
    attr_accessor :key, :value, :timestamp

    def initialize(pb_record)
      @key = pb_record.key
      @timestamp = pb_record.timestamp

      begin
        @value = JSON.parse(pb_record.value)
      rescue JSON::ParserError
        @value = pb_record.value
      end

      @value = @value.to_dot if is_value_hash?
    end

    def serialize
      Io::Meroxa::Funtime::Record.new(key: @key, value: @value.to_json, timestamp: @timestamp)
    end

    def get(key)
      if is_value_string? || is_value_array?
        @value
      elsif is_cdc_format?
        @value.send("payload.after.#{key}")
      else
        @value.send("payload.#{key}")
      end
    end

    def set(key, value)
      if !is_value_hash?
        @value = value
      else
        payload_key = is_cdc_format? ? "payload.after" : "payload"

        begin
          @value.send("#{payload_key}.#{key}")
        rescue NoMethodError
          schema = @value.send("schema.fields")
          new_schema_field = { field: key, optional: true, type: "string" }.to_dot

          if is_cdc_format?
            schema_fields = schema.find { |f| f.field == "after" }
            schema_fields.fields.unshift(new_schema_field)
          else
            schema.unshift(new_schema_field)
          end
        end

        @value.send("#{payload_key}.#{key}=", value)
      end
    end

    def unwrap!
      if is_cdc_format?
        payload = @value.send("payload")
        schema = @value.send("schema.fields")
        schema_fields = schema.find { |f| f.field == "after" }
        if !schema_fields.nil?
          schema_fields.delete("field")
          schema_fields.name = @value.send("schema.name")
          @value.send("schema=", schema_fields)
        end

        @value.send("payload=", payload.after)
      end
    end

    private

    def is_value_string?
      @value.is_a?(String)
    end

    def is_value_array?
      @value.is_a?(Array)
    end

    def is_value_hash?
      @value.is_a?(Hash)
    end

    def is_json_schema?
      is_value_hash? &&
      @value.has_key?("payload") &&
      @value.has_key?("schema")
    end

    def is_cdc_format?
      is_json_schema? &&
      @value.payload.has_key?("source")
    end

    def type_of_value(value)
      case value
      when String
        "string"
      when Integer
        "int32"
      when Float
        "float32"
      when true, false
        "boolean"
      else
        "unsupported"
      end
    end
  end

  class Records < SimpleDelegator
    def initialize(pb_records)
      records = pb_records.map { |r| Record.new(r) }
      __setobj__(records)
    end

    def unwrap!
      records = __getobj__
      records.each { |r| r.unwrap! }
    end
  end
end
