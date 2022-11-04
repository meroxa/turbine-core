# frozen_string_literal: true

require_relative 'lib/turbine/version'

Gem::Specification.new do |spec|
  spec.name = 'turbine'
  spec.version = Turbine::VERSION
  spec.authors = ['Ali Hamidi']
  spec.email = ['57750952+ahmeroxa@users.noreply.github.com']

  spec.summary = 'Turbine Data App framework Gem for Ruby.'
  spec.description = 'Turbine Framework Gem for building and deploying Turbine Data Apps on the Meroxa Data Platform.'
  spec.homepage = 'https://meroxa.com'
  spec.license = 'MIT'
  spec.required_ruby_version = '>= 2.6.0'

  spec.metadata['allowed_push_host'] = 'https://rubygems.org'

  spec.metadata['homepage_uri'] = spec.homepage
  spec.metadata['source_code_uri'] = 'https://github.com/turbine-core/lib/ruby/turbine'
  spec.metadata['changelog_uri'] = 'https://changelog.meroxa.com'

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`.split("\x0").reject do |f|
      (f == __FILE__) || f.match(%r{\A(?:(?:bin|test|spec|features)/|\.(?:git|travis|circleci)|appveyor)})
    end
  end
  spec.bindir = 'exe'
  spec.executables = spec.files.grep(%r{\Aexe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  # Uncomment to register a new dependency of your gem
  # spec.add_dependency "example-gem", "~> 1.0"
  spec.add_dependency 'grpc', "~> 1.48"

  # For more information and examples about making a new gem, check out our
  # guide at: https://bundler.io/guides/creating_gem.html
end
