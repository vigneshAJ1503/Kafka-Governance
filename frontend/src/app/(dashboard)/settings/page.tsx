'use client'

import { Card, CardHeader, CardTitle, CardContent, Input, Button } from '@/components/ui/core'
import { useAuthStore } from '@/store/useAuthStore'
import { useState } from 'react'

export default function SettingsPage() {
  const { user } = useAuthStore()
  const [theme, setTheme] = useState('dark')

  const toggleTheme = () => {
    const root = document.documentElement
    if (theme === 'dark') {
      root.classList.remove('dark')
      setTheme('light')
    } else {
      root.classList.add('dark')
      setTheme('dark')
    }
  }

  return (
    <div className="space-y-6 max-w-2xl">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Settings</h2>
        <p className="text-muted-foreground">Manage your account configurations and environment preferences.</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Profile Overview</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium text-muted-foreground">Username</label>
              <p className="font-medium text-foreground">{user?.username}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Current Role</label>
              <p className="font-medium text-foreground">{user?.role}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Appearance</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Dark Mode</p>
              <p className="text-sm text-muted-foreground">Toggle the UI theme style.</p>
            </div>
            <Button variant="outline" onClick={toggleTheme}>
              {theme === 'dark' ? 'Switch to Light' : 'Switch to Dark'}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Security Settings</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">New Password</label>
            <Input type="password" placeholder="••••••••" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Confirm Password</label>
            <Input type="password" placeholder="••••••••" />
          </div>
          <Button>Update Password</Button>
        </CardContent>
      </Card>

      {user?.role !== 'ADMIN' && (
         <Card>
           <CardHeader>
             <CardTitle className="text-lg text-primary">Elevate Access</CardTitle>
           </CardHeader>
           <CardContent className="flex flex-col gap-4 items-start">
              <p className="text-sm text-muted-foreground">
                You currently have <strong>{user?.role}</strong> access. Requesting ADMIN access will send a notification to existing cluster administrators for review.
              </p>
              <Button variant="secondary" onClick={() => alert('Access request submitted!')}>
                Request Admin Access
              </Button>
           </CardContent>
         </Card>
      )}

    </div>
  )
}
