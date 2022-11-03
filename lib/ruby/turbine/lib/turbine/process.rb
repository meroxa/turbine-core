# frozen_string_literal: true

require_relative '../proto/process/v1/service_services_pb'

module Turbine
  class Process < Io::Meroxa::Funtime::Function::Service
    def serve()
      function_addr = ENV['MEROXA_FUNCTION_ADDR'] ||= '0.0.0.0:50500'
      @grpc_server = GRPC::RpcServer.new
      @grpc_server.add_http2_port(function_addr, :this_port_is_insecure)
      @grpc_server.handle(self)
      puts "serving function #{self.class.name} on #{function_addr}"
      @grpc_server.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
    end

    def process(request, _call)
      processed_records = self.call(records: request.records)
      Io::Meroxa::Funtime::ProcessRecordResponse.new(records: processed_records.to_a)
    end
  end
end