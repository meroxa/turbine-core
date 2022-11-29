# frozen_string_literal: true
require 'service_services_pb'
require 'turbine_services_pb'

require "turbine_rb/collection_patch"
require "turbine_rb/version"
require "turbine_rb/client"

require 'optparse'
require 'fileutils'

module TurbineRb
  class Error < StandardError; end

  class << self
    attr_reader :app, :process_klass

    def register(app)
      @app = app
    end

    def register_fn(fn_klass)
      @process_klass = fn_klass
    end

    def serve
      process_function = @process_klass.new
      process_function_impl = ProcessImpl.new(process_function)
      function_addr = ENV['MEROXA_FUNCTION_ADDR'] ||= '0.0.0.0:50500'

      @grpc_server = GRPC::RpcServer.new
      @grpc_server.add_http2_port(function_addr, :this_port_is_insecure)
      @grpc_server.handle(process_function_impl)
      puts "serving function #{process_function.class.name} on #{function_addr}"
      @grpc_server.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
    end

    def run
      core_server = init_core_server()
      app = TurbineRb::Client::App.new(core_server)
      TurbineRb.app.call(app)
    end

    def record
      core_server = init_core_server()
      app = TurbineRb::Client::App.new(core_server, is_recording: true)
      TurbineRb.app.call(app)
    end

    def build
      docker_file = File.join(File.expand_path(File.dirname(__FILE__)) , "templates", "Dockerfile")
      dest_app = Dir.getwd
      FileUtils.cp(docker_file, dest_app)
    end

    private

    def init_core_server
      # TODO: figure out what the deal is with :this_channel_is_insecure
      core_server = TurbineCore::TurbineService::Stub.new(ENV["TURBINE_CORE_SERVER"], :this_channel_is_insecure)
      gitSHA = ARGV[0]

      req = TurbineCore::InitRequest.new(
        appName: app.class.name,
        configFilePath: Dir.getwd,
        language: :RUBY,
        gitSHA: gitSHA,
        turbineVersion: Gem.loaded_specs["turbine_rb"].version.version
      )

      core_server.init(req)
      core_server
    end
  end

  class ProcessImpl < Io::Meroxa::Funtime::Function::Service
    def initialize(process)
      @process = process
    end

    def process(request, _call)
      processed_records = @process.call(records: request.records)
      Io::Meroxa::Funtime::ProcessRecordResponse.new(records: processed_records.to_a)
    end
  end

  class Process
    def self.inherited(subclass)
      TurbineRb.register_fn(subclass)
      super
    end
  end
end
