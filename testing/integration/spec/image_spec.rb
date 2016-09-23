# encoding: utf-8
require 'spec_helper'

RSpec.describe "CLI fetch" do
  it "should remotely fetch a Docker container" do
    # TODO: This should do more stuff, probably
    cli.run!("fetch docker://nats")
    end
end
