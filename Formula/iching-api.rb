class IchingApi < Formula
  desc "A Go application for I Ching divination with Fiber-based HTTP API"
  homepage "https://github.com/peter7775/go-iching"
  url "https://github.com/peter7775/go-iching/releases/download/v#{version}/iching-api-darwin-amd64-v#{version}.tar.gz"
  sha256 "YOUR_SHA256_HERE" # Update this with actual SHA256 after release
  version "0.1.0"
  
  depends_on "go" => :build

  def install
    bin.install "iching-api-darwin-amd64" => "iching-api"
    pkgshare.install "static" => "web"
  end

  service do
    run [opt_bin/"iching-api"]
    keep_alive true
    log_path var/"log/iching-api.log"
    error_log_path var/"log/iching-api.error.log"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/iching-api --version", 1)
  end
end
