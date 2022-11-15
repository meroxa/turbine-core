# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: turbine/v1/turbine.proto

require 'google/protobuf'

require 'google/protobuf/empty_pb'
require 'google/protobuf/timestamp_pb'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("turbine/v1/turbine.proto", :syntax => :proto3) do
    add_message "turbine_core.InitRequest" do
      optional :appName, :string, 1
      optional :configFilePath, :string, 2
      optional :language, :enum, 3, "turbine_core.InitRequest.Language"
      optional :gitSHA, :string, 4
      optional :turbineVersion, :string, 5
    end
    add_enum "turbine_core.InitRequest.Language" do
      value :GOLANG, 0
      value :PYTHON, 1
      value :JAVASCRIPT, 2
      value :RUBY, 3
    end
    add_message "turbine_core.GetResourceRequest" do
      optional :name, :string, 1
    end
    add_message "turbine_core.Resource" do
      optional :uuid, :string, 1
      optional :name, :string, 2
      optional :type, :string, 3
      optional :direction, :enum, 4, "turbine_core.Resource.Direction"
    end
    add_enum "turbine_core.Resource.Direction" do
      value :DIRECTION_SOURCE, 0
      value :DIRECTION_DESTINATION, 1
    end
    add_message "turbine_core.Collection" do
      optional :name, :string, 1
      optional :stream, :string, 2
      repeated :records, :message, 3, "turbine_core.Record"
    end
    add_message "turbine_core.Record" do
      optional :key, :string, 1
      optional :value, :bytes, 2
      optional :timestamp, :message, 3, "google.protobuf.Timestamp"
    end
    add_message "turbine_core.ReadCollectionRequest" do
      optional :resource, :message, 1, "turbine_core.Resource"
      optional :collection, :string, 2
      optional :configs, :message, 3, "turbine_core.ResourceConfigs"
    end
    add_message "turbine_core.WriteCollectionRequest" do
      optional :resource, :message, 1, "turbine_core.Resource"
      optional :collection, :message, 2, "turbine_core.Collection"
      optional :targetCollection, :string, 3
      optional :configs, :message, 4, "turbine_core.ResourceConfigs"
    end
    add_message "turbine_core.ResourceConfigs" do
      repeated :resourceConfig, :message, 1, "turbine_core.ResourceConfig"
    end
    add_message "turbine_core.ResourceConfig" do
      optional :field, :string, 1
      optional :value, :string, 2
    end
    add_message "turbine_core.Process" do
      optional :name, :string, 1
    end
    add_message "turbine_core.ProcessCollectionRequest" do
      optional :process, :message, 1, "turbine_core.Process"
      optional :collection, :message, 2, "turbine_core.Collection"
    end
    add_message "turbine_core.Secret" do
      optional :name, :string, 1
      optional :value, :string, 2
    end
    add_message "turbine_core.ListFunctionsResponse" do
      repeated :functions, :string, 1
    end
    add_message "turbine_core.ResourceWithCollection" do
      optional :name, :string, 1
      optional :collection, :string, 2
      optional :direction, :enum, 3, "turbine_core.ResourceWithCollection.Direction"
    end
    add_enum "turbine_core.ResourceWithCollection.Direction" do
      value :DIRECTION_SOURCE, 0
      value :DIRECTION_DESTINATION, 1
    end
    add_message "turbine_core.ListResourcesResponse" do
      repeated :resources, :message, 1, "turbine_core.ResourceWithCollection"
    end
  end
end

module TurbineCore
  InitRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.InitRequest").msgclass
  InitRequest::Language = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.InitRequest.Language").enummodule
  GetResourceRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.GetResourceRequest").msgclass
  Resource = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Resource").msgclass
  Resource::Direction = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Resource.Direction").enummodule
  Collection = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Collection").msgclass
  Record = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Record").msgclass
  ReadCollectionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ReadCollectionRequest").msgclass
  WriteCollectionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.WriteCollectionRequest").msgclass
  ResourceConfigs = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ResourceConfigs").msgclass
  ResourceConfig = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ResourceConfig").msgclass
  Process = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Process").msgclass
  ProcessCollectionRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ProcessCollectionRequest").msgclass
  Secret = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.Secret").msgclass
  ListFunctionsResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ListFunctionsResponse").msgclass
  ResourceWithCollection = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ResourceWithCollection").msgclass
  ResourceWithCollection::Direction = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ResourceWithCollection.Direction").enummodule
  ListResourcesResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("turbine_core.ListResourcesResponse").msgclass
end