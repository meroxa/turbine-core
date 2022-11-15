# frozen_string_literal: true
module TurbineRb
  module Client
    class App
      attr_reader :core_server

      def initialize(grpc_server)
        @core_server = grpc_server
      end

      def resource(name:)
        req = TurbineCore::NameOrUUID.new(name:)
        res = @core_server.get_resource(req)
        Resource.new(res, self)
      end

      def process(records:, process:)
        records.pb_collection = process.call(records: records.pb_collection)
        records
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
  end
end
