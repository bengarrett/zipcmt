class Zipcmt < Formula
  desc "Zip Comment, the file comment viewer and extractor"
  homepage "https://github.com/bengarrett/zipcmt"
  url "https://github.com/bengarrett/zipcmt/archive/refs/tags/v1.4.7.tar.gz"
  sha256 "373dc71b5d6149b1b022830f0d66fd7c5c694b4763b4f0aa7969b6e4fb812ed1"
  version "1.4.7"
  license "LGPL-3.0-only"

  @commit = "41f2209be6f566f3d044a647bf61dde4cec16109"
  @build_date = "2026-05-01T15:58:46+10:00"

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
