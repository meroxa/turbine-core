# frozen_string_literal: true

require_relative "lib/turbine_rb/version"

Gem::Specification.new do |spec|
  spec.name = "turbine_rb"
  spec.version = TurbineRb::VERSION
  spec.authors = ["Ali Hamidi", "James Martinez"]
  spec.email = ["production@meroxa.io"]

  spec.summary = "Meroxa data application framework for Ruby"
  spec.description = "Turbine is a data application framework for building server-side applications that are event-driven, respond to data in real-time, and scale using cloud-native best practices"
  spec.homepage = "https://github.com/meroxa/turbine-core"
  spec.required_ruby_version = ">= 2.6.0"

  spec.metadata["homepage_uri"] = spec.homepage
  spec.metadata["source_code_uri"] = "https://github.com/meroxa/turbine-core"
  spec.metadata["changelog_uri"] = "https://github.com/meroxa/turbine-core"

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`.split("\x0").reject do |f|
      (f == __FILE__) || f.match(%r{\A(?:(?:bin|test|spec|features)/|\.(?:git|travis|circleci)|appveyor)})
    end
  end
  spec.bindir = "bin"
  spec.executables = ["turbine-function", "turbine-build", "turbine-record", "turbine-run"]
  spec.require_paths = ["lib"]

  # Uncomment to register a new dependency of your gem
  # spec.add_dependency "example-gem", "~> 1.0"
  spec.add_dependency "grpc", "~> 1.48"

  # For more information and examples about making a new gem, check out our
  # guide at: https://bundler.io/guides/creating_gem.html
end
