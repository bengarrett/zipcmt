class Zipcmt < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.8.tar.gz"
  sha256 "f143d4c6db95692f8e5f9b24d0b9fcbccc4656e0863f0d96b9b50d9d343c5ff9"
  version "1.4.8"
  license "LGPL-3.0-only"

  @commit = "ee925b236398543ab62b43560e2a0d0baf78a171"
  @build_date = "2026-08-07T22:45:40+10:00"

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
