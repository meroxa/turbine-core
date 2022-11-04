# frozen_string_literal: true

require_relative 'turbine/version'
require_relative 'turbine/process'
require_relative 'proto/turbine/v1/turbine_services_pb'

module Turbine
  class Error < StandardError; end

  class Runner
    attr_reader :core_server

    def initialize(app_name, gitSHA)
      turbine_server_addr = ENV['TURBINE_CORE_SERVER']

      # TODO: figure out what the deal is with :this_channel_is_insecure
      @core_server = TurbineCore::TurbineService::Stub.new(turbine_server_addr, :this_channel_is_insecure)
      req = TurbineCore::InitRequest.new(
        appName: app_name,
        configFilePath: Dir.getwd,
        language: :RUBY,
        gitSHA: gitSHA,
        turbineVersion: Gem.loaded_specs["turbine"].version.version
      )
      @core_server.init(req)
    end

    def resource(name:)
      req = TurbineCore::GetResourceRequest.new(name:)
      res = @core_server.get_resource(req)
      Resource.new(res, self)
    end

    def process(records:, process:)
      if records.instance_of?(Collection)
        unwraped_records = records.unwrap
      end
      p = TurbineCore::Process.new(
        name: process.class.name
      )
      
      req = TurbineCore::ProcessCollectionRequest.new(collection: unwraped_records, process: p)
      @core_server.add_process_to_collection(req)

      records.pb_collection = process.call(records: records.pb_collection)
      records
    end
  end

  class Resource
    attr_reader :pb_resource

    def initialize(res, app)
      @pb_resource = res
      @app = app
    end

    def records(collection:, configs: nil)
      req = TurbineCore::ReadCollectionRequest.new(resource: @pb_resource, collection:)
      req.configs = configs if configs
      @app.core_server.read_collection(req).wrap(@app) # wrap in Collection to enable chaining
    end

    def write(records:, collection:, configs: nil)
      if records.instance_of?(Collection) # it has been processed by a function, so unwrap back to gRPC collection
        records = records.unwrap
      end
      req = TurbineCore::WriteCollectionRequest.new(resource: @pb_resource, collection: records,
                                                targetCollection: collection)
      req.configs = configs if configs
      @app.core_server.write_collection_to_resource(req)
    end
  end

  class Collection
    attr_accessor :pb_collection, :pb_stream, :name

    def initialize(name, collection, stream, app)
      @name = name
      @pb_collection = collection
      @pb_stream = stream
      @app = app
    end

    def write_to(resource:, collection:, configs: nil)
      resource.write(records: self, collection:, configs:)
    end

    def process_with(process:)
      @app.process(records: self, process:)
    end

    def unwrap
      TurbineCore::Collection.new( # convert back to TurbineCore::Collection
        name:,
        records: pb_collection.to_a,
        stream: pb_stream
      )
    end
  end
end

TurbineCore::Collection.class_eval do
  def wrap(app)
    Turbine::Collection.new(
      name,
      records,
      stream,
      app
    )
  end
end

module Turbine
  attr_accessor :app

  def self.register(app)
    @app = app
  end

  def self.run
    gitSHA = ARGV[0]

    runner = Runner.new(@app.class.name, gitSHA)
    @app.call(runner)
  end
end
