# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: turbine_v2.proto for package 'turbine_core_v2'

require 'grpc'
require 'turbine_v2_pb'

module TurbineCoreV2
  module TurbineService
    class Service

      include ::GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'turbine_core_v2.TurbineService'

      rpc :Init, ::TurbineCoreV2::InitRequest, ::Google::Protobuf::Empty
      rpc :AddSource, ::TurbineCoreV2::AddSourceRequest, ::TurbineCoreV2::AddSourceResponse
      rpc :ReadRecords, ::TurbineCoreV2::ReadRecordsRequest, ::TurbineCoreV2::ReadRecordsResponse
      rpc :ProcessRecords, ::TurbineCoreV2::ProcessRecordsRequest, ::TurbineCoreV2::ProcessRecordsResponse
      rpc :AddDestination, ::TurbineCoreV2::AddDestinationRequest, ::TurbineCoreV2::AddDestinationResponse
      rpc :WriteRecords, ::TurbineCoreV2::WriteRecordsRequest, ::Google::Protobuf::Empty
      rpc :GetSpec, ::TurbineCoreV2::GetSpecRequest, ::TurbineCoreV2::GetSpecResponse
    end

    Stub = Service.rpc_stub_class
  end
end