RSpec.describe TurbineRb::Client::App do
  describe "#resource" do
    it "calls to grpc get_resource and returns a resource" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      stubs { |m| core_server.get_resource(m.is_a(TurbineCore::GetResourceRequest)) }.with { :resource }

      subject = TurbineRb::Client::App.new(core_server)
      result = subject.resource(name: "hey")

      expect(result.pb_resource).to eq(:resource)
      verify { |m| core_server.get_resource(m.that { |arg| arg.name == "hey" })}
    end
  end

  describe "#process" do
    let(:my_process) {
      Class.new(TurbineRb::Process) do
        def call(records:)
          records
        end
      end
    }
    it "calls the process function on the records" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      app = TurbineRb::Client::App.new(core_server)
      record = TurbineCore::Record.new(key: "1", value: "somebytes")
      records = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)

      result = app.process(records: records, process: my_process.new)

      expect(result.pb_collection.first.key).to eq("1")
    end
  end
end

RSpec.describe TurbineRb::Client::App::Resource do
  describe "#records" do
    it "calls to grpc read_collection and returns wrapped records" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      collection = Mocktail.of_next(TurbineCore::Collection)
      stubs { |m| core_server.read_collection(m.is_a(TurbineCore::ReadCollectionRequest)) }.with { TurbineCore::Collection.new }
      stubs { |m| collection.wrap(m.is_a(TurbineRb::Client::App)) }.with { :wrapped_collection }

      app = TurbineRb::Client::App.new(core_server)
      pb_resource = TurbineCore::Resource.new
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
      result = subject.records(collection: "hellocollection")

      expect(result).to eq(:wrapped_collection)
      verify {|m| core_server.read_collection(m.that { |arg| arg.collection == "hellocollection" }) }
      verify {|m| core_server.read_collection(m.that { |arg| arg.resource == pb_resource }) }
    end
  end

  describe "#write" do
    it "calls to grpc write_collection_to_resource" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      records = Mocktail.of(TurbineRb::Client::App::Collection)
      stubs { records.unwrap }.with { TurbineCore::Collection.new }

      app = TurbineRb::Client::App.new(core_server)
      pb_resource = TurbineCore::Resource.new
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
      subject.write(records: records, collection: "goodbyecollection")

      verify { |m| core_server.write_collection_to_resource(m.is_a(TurbineCore::WriteCollectionRequest)) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.resource == pb_resource}) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.targetCollection == "goodbyecollection"}) }
    end
  end
end

RSpec.describe TurbineRb::Client::App::Collection do
  describe "#write_to" do
    it "delegates to the resource write" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      resource = Mocktail.of(TurbineRb::Client::App::Resource)

      app = TurbineRb::Client::App.new(core_server)
      record = TurbineCore::Record.new(key: "1", value: "somebytes")

      subject = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)
      stubs { resource.write(records: subject, collection: "a_collection", configs: nil ) }.with { :write }

      result = subject.write_to(resource: resource, collection: "a_collection")
      expect(result).to eq(:write)
    end
  end

  describe "#process_with" do
    let(:my_process) {
      Class.new(TurbineRb::Process) do
        def call(records:)
          records
        end
      end
    }

    it "delegates to the app process" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      resource = Mocktail.of(TurbineRb::Client::App::Resource)
      app = Mocktail.of_next(TurbineRb::Client::App)

      record = TurbineCore::Record.new(key: "1", value: "somebytes")
      subject = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)
      process = my_process.new
      stubs { app.process(records: subject, process: process ) }.with { :process }

      result = subject.process_with(process: process)
      expect(result).to eq(:process)
    end
  end
end
