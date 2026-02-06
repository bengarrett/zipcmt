class Zipcmt < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.6.tar.gz"
  sha256 "d5558cd419c8d46bdc958064cb97f963d1ea793866414c025906ec15033512ed"
  version "1.4.6"
  license "LGPL-3.0-only"

  @commit = "4bb4c718fb9825efb22539b9311165837faacddc"
  @build_date = "2026-02-06T20:58:46+11:00"

  livecheck do
    url :stable
    strategy :github_latest
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version} -X main.commit=#{self.class.instance_variable_get('@commit')} -X main.date=#{self.class.instance_variable_get('@build_date')}")
  end

  test do
    assert_match "zipcmt", shell_output("#{bin}/zipcmt --version")
  end
end
