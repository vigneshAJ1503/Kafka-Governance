'use client'

import { useState } from 'react'
import { api } from '@/lib/api'
import { useAuthStore } from '@/store/useAuthStore'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { Card, CardHeader, CardTitle, CardContent, Input, Button } from '@/components/ui/core'

export default function SignupPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [role, setRole] = useState('USER') 
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const setCredentials = useAuthStore(state => state.setCredentials)
  const router = useRouter()

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const { data } = await api.post('/auth/signup', { username, password, role })
      setCredentials(data.token, data.user)
      router.push('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create account. Username might exist.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background/50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-3xl">Create Account</CardTitle>
          <p className="text-sm text-muted-foreground mt-2">Sign up for Kafka Governance access.</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSignup} className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Username</label>
              <Input value={username} onChange={e => setUsername(e.target.value)} required placeholder="engineer1" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Password</label>
              <Input type="password" value={password} onChange={e => setPassword(e.target.value)} required placeholder="••••••••" minLength={6} />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Requested Role</label>
              <select 
                value={role} 
                onChange={e => setRole(e.target.value)}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="USER">Standard User (Propose Topics)</option>
                <option value="ADMIN">Admin (Manage Policies & Approvals)</option>
              </select>
            </div>
            
            {error && <p className="text-sm text-destructive font-medium">{error}</p>}
            
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? 'Creating Account...' : 'Sign Up'}
            </Button>
            
            <div className="text-sm text-center mt-4">
              <span className="text-muted-foreground">Already have an account? </span>
              <Link href="/login" className="text-primary hover:underline font-medium">
                Sign In
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
