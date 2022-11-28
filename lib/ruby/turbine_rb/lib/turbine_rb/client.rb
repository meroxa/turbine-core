# frozen_string_literal: true

module TurbineRb
  module Client
    class App
      attr_reader :core_server

      def initialize(grpc_server, is_recording: false)
        @core_server = grpc_server
        @is_recording = is_recording
      end

      def resource(name:)
        req = TurbineCore::GetResourceRequest.new(name:)
        res = @core_server.get_resource(req)
        Resource.new(res, self)
      end

      def process(records:, process:)
        if records.instance_of?(Collection)
          unwrapped_records = records.unwrap
        end

        pr = TurbineCore::ProcessCollectionRequest::Process.new(
          name: process.class.name
        )

        req = TurbineCore::ProcessCollectionRequest.new(collection: unwrapped_records, process: pr)
        @core_server.add_process_to_collection(req)
        records.pb_collection = process.call(records: records.pb_collection) if not @is_recording

        records
      end

      # register_secrets accepts either a single string or an array of strings
      def register_secrets(secrets)
        [*secrets].map do |secret|
            req = TurbineCore::Secret.new(name: secret, value: ENV[secret])
            @core_server.register_secret(req)
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
          if configs
            pb_configs = configs.keys.map { |key| TurbineCore::Config.new(field: key, value: configs[key]) }
            req.configs = TurbineCore::Configs.new(config: pb_configs)
          end

          @app.core_server.read_collection(req).wrap(@app) # wrap in Collection to enable chaining
        end

        def write(records:, collection:, configs: nil)
          if records.instance_of?(Collection) # it has been processed by a function, so unwrap back to gRPC collection
            records = records.unwrap
          end

          req = TurbineCore::WriteCollectionRequest.new(resource: @pb_resource, sourceCollection: records,
                                                    targetCollection: collection)

          if configs
            pb_configs = configs.keys.map { |key| TurbineCore::Config.new(field: key, value: configs[key]) }
            req.configs = TurbineCore::Configs.new(config: pb_configs)
          end

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
  end
end
