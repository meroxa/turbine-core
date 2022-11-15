# frozen_string_literal: true

RSpec.describe TurbineRb do
  it "has a version number" do
    expect(TurbineRb::VERSION).not_to be nil
  end

  let(:app) { Class.new }
  let(:my_process) { Class.new(TurbineRb::Process) }

  describe ".register" do
    it "registers the app object" do
      stub_const("MyApp", app)
      my_app = MyApp.new
      TurbineRb.register(my_app)
      expect(TurbineRb.app).to eq(my_app)
    end
  end

  describe ".register_fn" do
    it "registers the function class" do
      stub_const("MyProcess", my_process)
      TurbineRb.register_fn(MyProcess)
      expect(TurbineRb.process_klass).to eq(MyProcess)
    end
  end

  describe ".serve" do
    it "serves the data app for funtime" do
      stub_const("MyProcess", my_process)
      TurbineRb.register_fn(MyProcess)
      grpc_server = Mocktail.of_next(GRPC::RpcServer)

      result = TurbineRb.serve
      verify { grpc_server.add_http2_port('0.0.0.0:50500', :this_port_is_insecure)}
      verify { |m| grpc_server.handle(m.is_a(TurbineRb::ProcessImpl)) }
    end
  end
end

RSpec.describe TurbineRb::ProcessImpl do
  describe "#process" do
    let(:my_process) {
      Class.new(TurbineRb::Process) do
        def call(records:)
          records
        end
      end
    }

    it "calls the function to process the records" do
      stub_const("MyProcess", my_process)
      record =  Io::Meroxa::Funtime::Record.new(key: "1", value: "somebytes")
      request = Io::Meroxa::Funtime::ProcessRecordRequest.new(records: [record])
      process = MyProcess.new
      subject = TurbineRb::ProcessImpl.new(process)
      result = subject.process(request, nil)

      expect(result).to be_instance_of(Io::Meroxa::Funtime::ProcessRecordResponse)
      expect(result.records.first.key).to eq("1")
    end
  end
end

RSpec.describe TurbineRb::Process do
  describe ".inherited" do
    let(:my_process) { Class.new(TurbineRb::Process) }
    it "calls to register the function" do
      stub_const("MyProcess", my_process)
      MyProcess.new
      expect(TurbineRb.process_klass).to eq(MyProcess)
    end
  end
end

