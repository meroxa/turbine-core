RSpec.describe TurbineRb::Record do
  describe "#serialize" do
    it "serializes the object to a proto record" do
      data = { key: "1",  value: { payload: { foo: "bar" }, schema: {} }.to_json, timestamp: Time.now.to_i}
      pb_record = Io::Meroxa::Funtime::Record.new(data)
      subject = TurbineRb::Record.new(pb_record)
      result = subject.serialize

      expect(result).to be_instance_of(Io::Meroxa::Funtime::Record)
    end
  end

  context "with cdc formatted json" do
    let(:data) do
      {
        key: "1",
        value: {
          payload: {
            after: {
              foo: "bar"
            },
            source: {}
          },
          schema: {
            fields: [
              {
                field: 'after',
                fields: [
                  { key: "foo", optional: false, type: "string" }
                ]
              }
            ],
            name: "schema_name"
          }
        }.to_json,
        timestamp: Time.now.to_i
      }
    end

    describe "#get" do
      it "returns the value at the key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        result = subject.get("foo")

        expect(result).to eq("bar")
      end
    end

    describe "#set" do
      it "sets the value at the key for an existing key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        subject.set("foo", "baz")
        result = subject.value.payload.after.foo

        expect(result).to eq("baz")
      end

      it "sets the value at the key and adds a schema entry for a new key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        subject.set("new_foo", "baz")
        value_result = subject.value.payload.after.new_foo
        schema_result = subject.value.schema
          .fields
          .find { |f| f.field == "after" }
          .fields
          .find { |f| f.field == "new_foo" }

        expect(value_result).to eq("baz")
        expect(schema_result.type).to eq("string")
      end
    end

    describe "unwrap!" do
      it "unwraps cdc formatted data into non cdc format" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        subject.unwrap!

        expect(subject.value.payload.foo).to eq("bar")
        expect(subject.value.schema.has_key?("field")).to be(false)
        expect(subject.value.schema.name).to eq("schema_name")
      end
    end
  end

  context "with non cdc formatted json" do
    let(:data) do
      {
        key: "1",
        value: {
          payload: {
            foo: "bar"
          },
          schema: {
            fields: [
              { key: "foo", optional: false, type: "string" }
            ]
          }
        }.to_json,
        timestamp: Time.now.to_i
      }
    end

    describe "#get" do
      it "returns the value at the key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        result = subject.get("foo")

        expect(result).to eq("bar")
      end
    end

    describe "#set" do
      it "sets the value at the key for an existing key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        subject.set("foo", "baz")
        result = subject.value.payload.foo

        expect(result).to eq("baz")
      end

      it "sets the value at the key and adds a schema entry for a new key" do
        pb_record = Io::Meroxa::Funtime::Record.new(data)
        subject = TurbineRb::Record.new(pb_record)
        subject.set("new_foo", "baz")
        value_result = subject.value.payload.new_foo
        schema_result = subject.value.schema.fields.find do |f|
          f.field == "new_foo"
        end

        expect(value_result).to eq("baz")
        expect(schema_result.type).to eq("string")
      end
    end
  end
end
