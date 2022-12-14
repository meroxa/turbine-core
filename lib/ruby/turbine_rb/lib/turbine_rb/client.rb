# frozen_string_literal: true

module TurbineRb
  module Client
    class MissingSecretError < StandardError; end

    class App
      attr_reader :core_server

      def initialize(grpc_server, recording:)
        @core_server = grpc_server
        @recording = recording
      end

      def recording?
        @recording
      end

      def resource(name:)
        req = TurbineCore::GetResourceRequest.new(name: name)
        res = @core_server.get_resource(req)
        Resource.new(res, self)
      end

      def process(records:, process:)
        new_records = records.dup
        pb_collection = core_server.add_process_to_collection(
          TurbineCore::ProcessCollectionRequest.new(
            collection: Collection.unwrap(new_records),
            process: TurbineCore::ProcessCollectionRequest::Process.new(name: process.class.name)
          )
        )
        new_records.tap do |r|
          r.pb_records = process_call(process: process, pb_collection: pb_collection)
          r.pb_stream = pb_collection.stream
        end
      end

      def process_call(process:, pb_collection:)
        return pb_collection.records if recording?

        process
          .call(records: TurbineRb::Records.new(pb_collection.records))
          .map(&:serialize_core_record)
      end

      # register_secrets accepts either a single string or an array of strings
      def register_secrets(secrets)
        [secrets].flatten.map do |secret|
          raise MissingSecretError, "secret #{secret} is not an environment variable" unless ENV.key?(secret)

          req = TurbineCore::Secret.new(name: secret, value: ENV[secret])
          core_server.register_secret(req)
        end
      end

      class Resource
        attr_reader :pb_resource, :app

        def initialize(res, app)
          @pb_resource = res
          @app = app
        end

        def records(collection:, configs: nil)
          req = TurbineCore::ReadCollectionRequest.new(resource: @pb_resource, collection: collection)
          if configs
            pb_configs = configs.keys.map { |key| TurbineCore::Config.new(field: key, value: configs[key]) }
            req.configs = TurbineCore::Configs.new(config: pb_configs)
          end

          app.core_server
             .read_collection(req)
             .wrap(app) # wrap in Collection to enable chaining
        end

        def write(records:, collection:, configs: nil)
          new_records = records.dup
          if records.instance_of?(Collection) # it has been processed by a function, so unwrap back to gRPC collection
            new_records = records.unwrap
          end

          req = TurbineCore::WriteCollectionRequest.new(
            resource: @pb_resource,
            sourceCollection: new_records,
            targetCollection: collection
          )

          if configs
            pb_configs = configs.keys.map { |key| TurbineCore::Config.new(field: key, value: configs[key]) }
            req.configs = TurbineCore::Configs.new(config: pb_configs)
          end

          app.core_server.write_collection_to_resource(req)
        end
      end

      class Collection
        attr_accessor :pb_records, :pb_stream, :name, :app

        def self.unwrap(collection)
          return collection.unwrap if collection.instance_of?(Collection)

          collection
        end

        def initialize(name, records, stream, app)
          @name = name
          @pb_records = records
          @pb_stream = stream
          @app = app
        end

        def write_to(resource:, collection:, configs: nil)
          resource.write(records: self, collection: collection, configs: configs)
        end

        def process_with(process:)
          app.process(records: self, process: process)
        end

        def records
          return pb_records.to_a if pb_records.respond_to?(:to_a)

          [pb_records].compact.flatten
        end

        def unwrap
          TurbineCore::Collection.new( # convert back to TurbineCore::Collection
            name: name,
            records: records,
            stream: pb_stream
          )
        end
      end
    end
  end
end
