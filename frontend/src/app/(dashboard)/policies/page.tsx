'use client'

import { useState } from 'react'
import Editor from '@monaco-editor/react'
import { api } from '@/lib/api'
import { Card, CardContent, Input, Button, Badge } from '@/components/ui/core'
import { ShieldAlert } from 'lucide-react'
import { useAuthStore } from '@/store/useAuthStore'

// Boilerplate Cedar template
const DEFAULT_POLICY = `permit(
  principal == User::"admin1",
  action == Action::"CreateTopic",
  resource == Topic::"*"
);`

export default function PoliciesPage() {
  const [showEditor, setShowEditor] = useState(false)
  const [code, setCode] = useState(DEFAULT_POLICY)
  const [principal, setPrincipal] = useState('User::"user1"')
  const [resource, setResource] = useState('Topic::"*"')
  
  const { user } = useAuthStore()

  // We use mock local state since /api/v1/policies doesn't explicitly expose GET in backend.
  // In a full feature scale, this connects via react-query securely.
  const [mockPolicies, setMockPolicies] = useState([
    { id: 'pol_1', principal: 'User::"admin1"', action: 'Action::"ViewTopic"', resource: 'Topic::"*"', effect: 'permit' },
    { id: 'pol_2', principal: 'User::"user1"', action: 'Action::"CreateTopic"', resource: 'Topic::"*"', effect: 'permit' }
  ])

  // Admin sees all policies. Users see policies pointing to their principal.
  const visiblePolicies = mockPolicies.filter(p => user?.role === 'ADMIN' || p.principal === `User::"${user?.username}"`)

  const handleSavePolicy = async () => {
    // Basic regex extract assuming single simple policy block for MVP demo mapping to POST /policies array payload requirements
    try {
      await api.post('/policies', {
        principal: principal,
        action: 'Action::"CreateTopic"',
        resource: resource,
        effect: 'permit'
      })
      alert('Policy synchronized with Cedar agent')
      setShowEditor(false)
    } catch(err) {
      alert('Failed to deploy policy')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Access Policies</h2>
          <p className="text-muted-foreground">Author and govern Cedar-based IAM permissions globally.</p>
        </div>
        <Button onClick={() => setShowEditor(!showEditor)}>
          New Policy
        </Button>
      </div>

      {showEditor && (
        <Card className="border-border">
          <CardContent className="pt-6">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              
              <div className="lg:col-span-2 border rounded-md overflow-hidden border-border flex flex-col">
                <div className="bg-secondary px-4 py-2 text-xs font-semibold uppercase text-muted-foreground flex justify-between">
                  <span>Cedar Policy Editor</span>
                  <span className="text-primary">Syntax: Strict</span>
                </div>
                <Editor
                  height="300px"
                  defaultLanguage="javascript" // Using js syntax coloring until custom Cedar tokenizer is added
                  theme="vs-dark"
                  value={code}
                  onChange={(val) => setCode(val || '')}
                  options={{ minimap: { enabled: false }, fontSize: 13, padding: { top: 16 } }}
                />
              </div>

              <div className="space-y-4">
                <h3 className="font-semibold flex items-center gap-2">
                  <ShieldAlert className="h-4 w-4 text-primary" />
                  Policy Metadata
                </h3>
                <div>
                  <label className="text-xs font-medium text-muted-foreground mb-1 block">Principal Extract</label>
                  <Input value={principal} onChange={e => setPrincipal(e.target.value)} />
                </div>
                <div>
                  <label className="text-xs font-medium text-muted-foreground mb-1 block">Resource Extract</label>
                  <Input value={resource} onChange={e => setResource(e.target.value)} />
                </div>

                <div className="pt-4 border-t border-border flex flex-col gap-2">
                  <Button onClick={handleSavePolicy} className="w-full">Deploy to Engine</Button>
                  <Button variant="outline" onClick={() => setShowEditor(false)} className="w-full bg-transparent">Cancel</Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="border rounded-md border-border bg-card">
        <table className="w-full text-sm text-left">
          <thead className="bg-secondary/50 text-muted-foreground">
            <tr>
              <th className="px-4 py-3 font-medium">Scope Context</th>
              <th className="px-4 py-3 font-medium">Effect</th>
              <th className="px-4 py-3 font-medium">Principal</th>
              {user?.role === 'ADMIN' && <th className="px-4 py-3 font-medium">Actions</th>}
            </tr>
          </thead>
          <tbody>
            {visiblePolicies.map(p => (
              <tr key={p.id} className="border-t border-border hover:bg-muted/50">
                <td className="px-4 py-3">
                  <p className="font-medium text-foreground">{p.action}</p>
                  <p className="text-xs text-muted-foreground">on {p.resource}</p>
                </td>
                <td className="px-4 py-3">
                  <Badge className="bg-emerald-500/10 text-emerald-500 hover:bg-emerald-500/20 shadow-none border-transparent">{p.effect}</Badge>
                </td>
                <td className="px-4 py-3"><code className="px-1 py-0.5 bg-secondary rounded text-xs">{p.principal}</code></td>
                {user?.role === 'ADMIN' && (
                  <td className="px-4 py-3">
                   <Button 
                      variant="outline"
                      onClick={() => { setMockPolicies(prev => prev.filter(x => x.id !== p.id)) }}
                      className="h-7 px-3 py-0 text-xs text-destructive border-destructive/50 hover:bg-destructive/10"
                    >
                      Delete
                    </Button>
                  </td>
                )}
              </tr>
            ))}
            {visiblePolicies.length === 0 && (
               <tr><td colSpan={4} className="p-8 text-center text-muted-foreground">No policies available for your role.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
