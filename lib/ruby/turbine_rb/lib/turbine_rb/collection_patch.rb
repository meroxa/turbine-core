# frozen_string_literal: true

TurbineCore::Collection.class_eval do
  def wrap(app)
    TurbineRb::Client::App::Collection.new(
      name,
      records,
      stream,
      app
    )
  end
end
