RSpec.describe TurbineRb::Client::App do
  describe "#resource" do
    it "calls to grpc get_resource and returns a resource" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      stubs { |m| core_server.get_resource(m.is_a(TurbineCore::GetResourceRequest)) }.with { :resource }

      subject = TurbineRb::Client::App.new(core_server)
      result = subject.resource(name: "hey")

      expect(result.pb_resource).to eq(:resource)
      verify { |m| core_server.get_resource(m.that { |arg| arg.name == "hey" }) }
    end
  end

  describe "#process" do
    let(:my_process) do
      Class.new(TurbineRb::Process) do
        def call(records:)
          records.first.value = "changedbytes"
          records
        end
      end
    end

    let(:core_server) { Mocktail.of(TurbineCore::TurbineService::Stub) }
    let(:record) { TurbineCore::Record.new(key: "1", value: "somebytes") }

    it "calls the process function on the records in run mode" do
      app = TurbineRb::Client::App.new(core_server)
      records = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)
      result = app.process(records: records, process: my_process.new)

      expect(result.pb_collection.first.key).to eq("1")
      expect(result.pb_collection.first.value).to eq("changedbytes")
    end

    it "doesnt call the process function on the records in record mode" do
      app = TurbineRb::Client::App.new(core_server, is_recording: true)
      records = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)
      result = app.process(records: records, process: my_process.new)

      expect(result.pb_collection.first.key).to eq("1")
      expect(result.pb_collection.first.value).to eq("somebytes")
    end
  end

  describe "#register_secret" do
    let(:secrets) do
      [{ name: "ENV_VAR", value: "value" }, { name: "ENV_VAR_2", value: "value_2" }]
    end

    let(:core_server) do
      Mocktail.of(TurbineCore::TurbineService::Stub)
    end

    let(:app) do
      TurbineRb::Client::App.new(core_server)
    end

    before do
      secrets.each do |s|
        ENV[s[:name]] = s[:value]
      end

      stubs { |m| core_server.register_secret(m.is_a(TurbineCore::Secret)) }.with { TurbineCore::Secret.new }
    end

    after do
      secrets.each do |s|
        ENV.delete(s[:name])
      end
    end

    it "raises an error when secret is missing from env" do
      expect do
        app.register_secrets("FOOBAR")
      end.to raise_error(
        TurbineRb::Client::MissingSecretError,
        /FOOBAR is not an environment variable/
      )
    end

    it "calls to grpc register_secret using a single secret" do
      user_secret = secrets[0][:name]
      app.register_secrets(user_secret)

      verify(times: 1) do |m|
        core_server.register_secret(m.that do |arg|
          arg.name == secrets[0][:name] && arg.value == secrets[0][:value]
        end)
      end
    end

    it "calls to grpc register_secret using an array of secrets" do
      user_secrets = [secrets[0][:name], secrets[1][:name]]
      app.register_secrets(user_secrets)

      2.times do |i|
        verify(times: 1) do |m|
          core_server.register_secret(m.that do |arg|
            arg.name == secrets[i][:name] && arg.value == secrets[i][:value]
          end)
        end
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
      stubs do |m|
        core_server.read_collection(m.is_a(TurbineCore::ReadCollectionRequest))
      end.with { TurbineCore::Collection.new }
      stubs { |m| collection.wrap(m.is_a(TurbineRb::Client::App)) }.with { :wrapped_collection }
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
    end

    it "calls to grpc read_collection and returns wrapped records" do
      result = subject.records(collection: "hellocollection")

      expect(result).to eq(:wrapped_collection)
      verify { |m| core_server.read_collection(m.that { |arg| arg.collection == "hellocollection" }) }
      verify { |m| core_server.read_collection(m.that { |arg| arg.resource == pb_resource }) }
    end

    it "sets configuration when configs arg is passed" do
      result = subject.records(collection: "hellocollection", configs: { "some.key" => "some.value" })

      verify { |m| core_server.read_collection(m.that { |arg| arg.configs.config.first.field == "some.key" }) }
      verify { |m| core_server.read_collection(m.that { |arg| arg.configs.config.first.value == "some.value" }) }
    end
  end

  describe "#write" do
    let(:core_server) { Mocktail.of(TurbineCore::TurbineService::Stub) }
    let(:records) { Mocktail.of(TurbineRb::Client::App::Collection) }
    let(:pb_resource) { TurbineCore::Resource.new }
    let(:app) { TurbineRb::Client::App.new(core_server) }
    let(:subject) do
      stubs { records.unwrap }.with { TurbineCore::Collection.new }
      subject = TurbineRb::Client::App::Resource.new(pb_resource, app)
    end

    it "calls to grpc write_collection_to_resource" do
      subject.write(records: records, collection: "goodbyecollection")

      verify { |m| core_server.write_collection_to_resource(m.is_a(TurbineCore::WriteCollectionRequest)) }
      verify { |m| core_server.write_collection_to_resource(m.that { |arg| arg.resource == pb_resource }) }
      verify do |m|
        core_server.write_collection_to_resource(m.that do |arg|
                                                   arg.targetCollection == "goodbyecollection"
                                                 end)
      end
    end

    it "sets configuration when configs arg is passed" do
      subject.write(records: records, collection: "goodbyecollection", configs: { "some.key" => "some.value" })
      verify do |m|
        core_server.write_collection_to_resource(m.that do |arg|
                                                   arg.configs.config.first.field == "some.key"
                                                 end)
      end
      verify do |m|
        core_server.write_collection_to_resource(m.that do |arg|
                                                   arg.configs.config.first.value == "some.value"
                                                 end)
      end
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
      stubs { resource.write(records: subject, collection: "a_collection", configs: nil) }.with { :write }

      result = subject.write_to(resource: resource, collection: "a_collection")
      expect(result).to eq(:write)
    end
  end

  describe "#process_with" do
    let(:my_process) do
      Class.new(TurbineRb::Process) do
        def call(records:)
          records
        end
      end
    end

    it "delegates to the app process" do
      core_server = Mocktail.of(TurbineCore::TurbineService::Stub)
      resource = Mocktail.of(TurbineRb::Client::App::Resource)
      app = Mocktail.of_next(TurbineRb::Client::App)

      record = TurbineCore::Record.new(key: "1", value: "somebytes")
      subject = TurbineRb::Client::App::Collection.new("a_name", [record], "a_stream", app)
      process = my_process.new
      stubs { app.process(records: subject, process: process) }.with { :process }

      result = subject.process_with(process: process)
      expect(result).to eq(:process)
    end
  end
end
