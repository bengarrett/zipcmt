class Namzd < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.5.tar.gz"
  sha256 "17446067087fafdbfdcf207bd3d1b52d1f872e736f6a7c05b2666018b08bb582"
  version "1.4.5"
  license "LGPL-3.0-only"

  @commit = "b647d43e354131b5c17f7d6fc3696f2833f990a4"
  @build_date = "2026-02-06T12:11:20+11:00"

  livecheck do
    url :stable
    strategy :github_latest
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version} -X main.commit=#{self.class.instance_variable_get('@commit')} -X main.date=#{self.class.instance_variable_get('@build_date')}")
  end

  test do
    assert_match "namzd", shell_output("#{bin}/namzd --version")
  end
end
