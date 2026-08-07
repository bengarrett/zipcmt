class Zipcmt < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.8.tar.gz"
  sha256 "c52dc03bb494814bef1229028895478e64ecbec88cee220dc5c26427092f7b5e"
  version "1.4.8"
  license "LGPL-3.0-only"

  @commit = "e70adc816d4d47308e8e57eec676e6a02590aabc"
  @build_date = "2026-08-07T22:39:53+10:00"

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
