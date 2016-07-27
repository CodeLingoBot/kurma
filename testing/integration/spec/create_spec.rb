# encoding: utf-8
require 'spec_helper'

RSpec.describe "CLI create" do
  it "should launch an AppC container" do
    initial_pods_count = api.list_pods["pods"].size

    output = cli.run!("create coreos.com/etcd:v2.2.5")
    uuid = output.scan(/Launched pod ([\w-]+)/).flatten.first
    expect(uuid).not_to be_nil
    @cleanup << "stop #{uuid}"

    resp = api.list_pods
    expect(resp["pods"].size).to eq(initial_pods_count+1)

    output = cli.run!("stop #{uuid}")
    expect(output).to include("Destroyed pod")

    resp = api.list_pods
    expect(resp["pods"].size).to eq(initial_pods_count)
  end

  it "should launch a Docker container" do
    initial_pods_count = api.list_pods["pods"].size

    output = cli.run!("create docker://nats")
    uuid = output.scan(/Launched pod ([\w-]+)/).flatten.first
    expect(uuid).not_to be_nil
    @cleanup << "stop #{uuid}"

    resp = api.list_pods
    expect(resp["pods"].size).to eq(initial_pods_count+1)
  end

  describe "enter a container" do
    it "should run a command" do
      output = cli.run!("create docker://busybox --name busybox --net=host /bin/sleep 60")
      uuid = output.scan(/Launched pod ([\w-]+)/).flatten.first
      expect(uuid).not_to be_nil
      @cleanup << "stop #{uuid}"

      output = cli.run!("enter #{uuid} busybox", "whoami", "exit")
      output.gsub!("\r", "") # trim carriage returns
      expect(output).to match("whoami\nroot")
    end

    it "should return the status code on exit" do
      output = cli.run!("create docker://busybox --name busybox --net=host /bin/sleep 60")
      uuid = output.scan(/Launched pod ([\w-]+)/).flatten.first
      expect(uuid).not_to be_nil
      @cleanup << "stop #{uuid}"

      _, _, status = cli.run("enter #{uuid} busybox", "exit 45")
      expect(status.exitstatus).to be(45)
    end

    it "should fork off a process" do
      output = cli.run!("create docker://busybox --name busybox --net=host /bin/sleep 60")
      uuid = output.scan(/Launched pod ([\w-]+)/).flatten.first
      expect(uuid).not_to be_nil
      @cleanup << "stop #{uuid}"

      time_before = Time.now.to_i
      output = cli.run!("enter #{uuid} busybox", "sleep 5 &", "echo forked", "exit")
      time_after = Time.now.to_i

      expect(time_after - time_before).to be < 5

      output.gsub!("\r", "") # trim carriage returns
      expect(output).to match("forked")

      # wait 7 seconds, the pod should still be alive
      sleep 7

      # get the pod, validate its status
      pod = api.get_pod(uuid)
      expect(pod["pod"]["state"]).to eq("RUNNING")
    end
  end
end
