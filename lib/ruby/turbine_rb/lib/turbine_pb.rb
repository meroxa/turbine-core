# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: turbine.proto

require 'google/protobuf'

require 'google/protobuf/empty_pb'
require 'google/protobuf/timestamp_pb'
require 'google/protobuf/wrappers_pb'
require 'validate/validate_pb'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("turbine.proto", :syntax => :proto3) do
    add_message "turbine_core.ReadFromSourceRequest" do
      optional :pluginName, :string, 1
      optional :direction, :enum, 2, "turbine_core.ReadFromSourceRequest.Direction"
      map :configuration, :string, :string, 3
    end
    add_enum "turbine_core.ReadFromSourceRequest.Direction" do
      value :SOURCE, 0
      value :DESTINATION, 1
    end
    add_message "turbine_core.WriteToDestinationRequest" do
      optional :pluginName, :string, 1
      map :configuration, :string, :string, 2
      optional :records, :message, 3, "turbine_core.RecordsCollection"
    end
    add_message "turbine_core.InitRequest" do
      optional :appName, :string, 1
      optional :configFilePath, :string, 2
      optional :language, :enum, 3, "turbine_core.Language"
      optional :gitSHA, :string, 4
      optional :turbineVersion, :string, 5
    end
    add_message "turbine_core.GetResourceRequest" do
      optional :name, :string, 1
    end
    add_message "turbine_core.Resource" do
      optional :name, :string, 1
      optional :source, :bool, 2
      optional :destination, :bool, 3
      optional :collection, :string, 4
    end
    add_message "turbine_core.RecordsCollection" do
      repeated :records, :message, 1, "turbine_core.Record"
    end
    add_message "turbine_core.Record" do
      optional :key, :string, 1
      optional :value, :bytes, 2
      optional :timestamp, :message, 3, "google.protobuf.Timestamp"
    end
    add_message "turbine_core.ReadCollectionRequest" do
      optional :resource, :message, 1, "turbine_core.Resource"
      optional :collection, :string, 2
      optional :configs, :message, 3, "turbine_core.Configs"
    end
    add_message "turbine_core.WriteCollectionRequest" do
      optional :resource, :message, 1, "turbine_core.Resource"
      optional :sourceCollection, :message, 2, "turbine_core.RecordsCollection"
      optional :targetCollection, :string, 3
      optional :configs, :message, 4, "turbine_core.Configs"
    end
    add_message "turbine_core.Configs" do
      repeated :config, :message, 1, "turbine_core.Config"
    end
    add_message "turbine_core.Config" do
      optional :field, :string, 1
      optional :value, :string, 2
    end
    add_message "turbine_core.ProcessRecordsRequest" do
      optional :process, :message, 1, "turbine_core.ProcessRecordsRequest.Process"
      repeated :records, :message, 2, "turbine_core.Record"
    end
    add_message "turbine_core.ProcessRecordsRequest.Process" do
      optional :name, :string, 1
    end
    add_message "turbine_core.Secret" do
      optional :name, :string, 1
      optional :value, :string, 2
    end
    add_message "turbine_core.ListResourcesResponse" do
      repeated :resources, :message, 1, "turbine_core.Resource"
    end
    add_message "turbine_core.GetSpecRequest" do
      optional :image, :string, 1
    end
    add_message "turbine_core.GetSpecResponse" do
      optional :spec, :bytes, 1
    end
    add_enum "turbine_core.Language" do
      value :GOLANG, 0
      value :PYTHON, 1
      value :JAVASCRIPT, 2
      value :RUBY, 3
      value :JAVA, 4
    end
  end
end

module TurbineCore
  ReadFromSourceRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ReadFromSourceRequest").msgclass
  ReadFromSourceRequest::Direction = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ReadFromSourceRequest.Direction").enummodule
  WriteToDestinationRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.WriteToDestinationRequest").msgclass
  InitRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.InitRequest").msgclass
  GetResourceRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.GetResourceRequest").msgclass
  Resource = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Resource").msgclass
  RecordsCollection = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.RecordsCollection").msgclass
  Record = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Record").msgclass
  ReadCollectionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ReadCollectionRequest").msgclass
  WriteCollectionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.WriteCollectionRequest").msgclass
  Configs = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Configs").msgclass
  Config = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Config").msgclass
  ProcessRecordsRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ProcessRecordsRequest").msgclass
  ProcessRecordsRequest::Process = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ProcessRecordsRequest.Process").msgclass
  Secret = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Secret").msgclass
  ListResourcesResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ListResourcesResponse").msgclass
  GetSpecRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.GetSpecRequest").msgclass
  GetSpecResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.GetSpecResponse").msgclass
  Language = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Language").enummodule
end
