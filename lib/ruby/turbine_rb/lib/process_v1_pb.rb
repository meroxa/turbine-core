# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: process_v1.proto

require 'google/protobuf'

require 'google/protobuf/struct_pb'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("process_v1.proto", :syntax => :proto3) do
    add_message "io.meroxa.funtime.ProcessRecordRequest" do
      repeated :records, :message, 1, "io.meroxa.funtime.Record"
    end
    add_message "io.meroxa.funtime.ProcessRecordResponse" do
      repeated :records, :message, 1, "io.meroxa.funtime.Record"
    end
    add_message "io.meroxa.funtime.Record" do
      optional :key, :string, 1
      optional :value, :string, 2
      optional :timestamp, :int64, 3
      optional :structured, :message, 4, "google.protobuf.Struct"
    end
  end
end

module Io
  module Meroxa
    module Funtime
      ProcessRecordRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("io.meroxa.funtime.ProcessRecordRequest").msgclass
      ProcessRecordResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("io.meroxa.funtime.ProcessRecordResponse").msgclass
      Record = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("io.meroxa.funtime.Record").msgclass
    end
  end
end
