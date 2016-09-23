# encoding: utf-8
require 'spec_helper'

RSpec.describe "CLI fetch" do
  it "should remotely fetch a Docker container" do
    initial_images_count = api.list_images["images"].size

    cli.run!("fetch docker://nats")

    resp = api.list_images
    expect(resp["images"].size).to eq(initial_images_count+1)
  end

  it "should remotely fetch an appc container" do
    initial_images_count = api.list_images["images"].size

    cli.run!("fetch coreos.com/etcd:v2.2.5")

    resp = api.list_images
    expect(resp["images"].size).to eq(initial_images_count+1)
  end
end
