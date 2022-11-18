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

  describe "#register_secret" do
    let(:secrets) {
      [{ name: "ENV_VAR", value: "value"}, { name: "ENV_VAR_2", value: "value_2"}]
    }

    let(:core_server) {
      Mocktail.of(TurbineCore::TurbineService::Stub)
    }

    let(:app) {
      TurbineRb::Client::App.new(core_server)
    }

    before(:each) do
      secrets.each do |s|
        ENV[s[:name]] = s[:value]
      end

      stubs { |m| core_server.register_secret(m.is_a(TurbineCore::Secret)) }.with { TurbineCore::Secret.new }
    end

    after(:each) do
      secrets.each do |s|
        ENV.delete(s[:name])
      end
    end

    it "calls to grpc register_secret using a single secret" do
      user_secret = secrets[0][:name]
      app.register_secrets(user_secret)

      verify(times: 1) { |m|
        core_server.register_secret(m.that { |arg|
          arg.name == secrets[0][:name] && arg.value == secrets[0][:value]
        })
      }
    end

    it "calls to grpc register_secret using an array of secrets" do
      user_secrets = [secrets[0][:name], secrets[1][:name]]
      app.register_secrets(user_secrets)

      2.times do |i|
        verify(times: 1) { |m|
          core_server.register_secret(m.that { |arg|
            arg.name == secrets[i][:name] && arg.value == secrets[i][:value]
          })
        }
      end
    end
  end
end

RSpec.describe TurbineRb::Client::App::Resource do
  describe "#records" do
    let(:core_server) { Mocktail.of(TurbineCore::TurbineService::Stub) }
    let(:collection) { Mocktail.of_next(TurbineCore::Collection) }
    let(:pb_resource) { TurbineCore::Resource.new }
    let(:app) { TurbineRb::Client::App.new(core_server) }
    let(:subject) do
      stubs { |m| core_server.read_collection(m.is_a(TurbineCore::ReadCollectionRequest)) }.with { TurbineCore::Collection.new }
      stubs { |m| collection.wrap(m.is_a(TurbineRb::Client::App)) }.with { :wrapped_collection }
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
    end

    it "calls to grpc read_collection and returns wrapped records" do
      result = subject.records(collection: "hellocollection")

      expect(result).to eq(:wrapped_collection)
      verify {|m| core_server.read_collection(m.that { |arg| arg.collection == "hellocollection" }) }
      verify {|m| core_server.read_collection(m.that { |arg| arg.resource == pb_resource }) }
    end

    it "sets configuration when configs arg is passed" do
      result = subject.records(collection: "hellocollection", configs: { "some.key" => "some.value" })

      verify { |m| core_server.read_collection(m.that { |arg| arg.configs.config.first.field == "some.key"}) }
      verify { |m| core_server.read_collection(m.that { |arg| arg.configs.config.first.value == "some.value"}) }
    end
  end

  describe "#write" do
    let(:core_server) { Mocktail.of(TurbineCore::TurbineService::Stub) }
    let(:records) {  Mocktail.of(TurbineRb::Client::App::Collection) }
    let(:pb_resource) { TurbineCore::Resource.new }
    let(:app) { TurbineRb::Client::App.new(core_server) }
    let(:subject) do
      stubs { records.unwrap }.with { TurbineCore::Collection.new }
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
    end

    it "calls to grpc write_collection_to_resource" do
      subject.write(records: records, collection: "goodbyecollection")

      verify { |m| core_server.write_collection_to_resource(m.is_a(TurbineCore::WriteCollectionRequest)) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.resource == pb_resource}) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.targetCollection == "goodbyecollection"}) }
    end

    it "sets configuration when configs arg is passed" do
      subject.write(records: records, collection: "goodbyecollection", configs: { "some.key" => "some.value" })
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.configs.config.first.field == "some.key"}) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.configs.config.first.value == "some.value"}) }
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
