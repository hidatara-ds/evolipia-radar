/** @type {import('next').NextConfig} */
const nextConfig = {
  // Rewrites to avoid conflict with Go API during dev
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: '/api/:path*',
      },
    ]
  },
}

module.exports = nextConfig
