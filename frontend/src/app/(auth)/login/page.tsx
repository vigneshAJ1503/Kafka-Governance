'use client'

import { useState } from 'react'
import { api } from '@/lib/api'
import { useAuthStore } from '@/store/useAuthStore'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Card, CardHeader, CardTitle, CardContent, Input, Button } from '@/components/ui/core'

export default function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const setCredentials = useAuthStore(state => state.setCredentials)
  const router = useRouter()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const { data } = await api.post('/auth/signin', { username, password })
      setCredentials(data.token, data.user)
      router.push('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Invalid credentials')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background/50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-3xl">Kafka Governance</CardTitle>
          <p className="text-sm text-muted-foreground mt-2">Sign in to manage policies and topics.</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleLogin} className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Username</label>
              <Input value={username} onChange={e => setUsername(e.target.value)} required placeholder="admin1" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Password</label>
              <Input type="password" value={password} onChange={e => setPassword(e.target.value)} required placeholder="••••••••" />
            </div>
            {error && <p className="text-sm text-destructive font-medium">{error}</p>}
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Authenticating...' : 'Sign In'}
            </Button>
            
            <div className="text-sm text-center mt-4">
              <span className="text-muted-foreground">Don't have an account? </span>
              <Link href="/signup" className="text-primary hover:underline font-medium">
                Sign Up
              </Link>
            </div>
            
            <div className="text-xs text-muted-foreground text-center mt-4 border-t border-border pt-4">
              Demo instances use role-based routing (e.g. Try admin1 / password123).
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
