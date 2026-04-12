/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://kafka-governance:8080/api/:path*', // Proxy to Backend
      },
    ]
  },
}
module.exports = nextConfig
