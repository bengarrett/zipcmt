class Zipcmt < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.9.tar.gz"
  sha256 "4d3e241a9c37c48a8230b9ea54493ddfffe26368114a1935609a4be653fed7e7"
  version "1.4.9"
  license "LGPL-3.0-only"

  @commit = "a59e1e804f80c51ba441d6555fd0a5ee4cbbd81a"
  @build_date = "2026-10-28T09:40:55+11:00"

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
