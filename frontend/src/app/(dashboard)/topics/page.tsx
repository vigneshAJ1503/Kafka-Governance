'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { Card, CardContent, Input, Button, Badge } from '@/components/ui/core'
import { useAuthStore } from '@/store/useAuthStore'
import { Plus, ShieldAlert, CheckCircle2 } from 'lucide-react'

export default function TopicsPage() {
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [showCreate, setShowCreate] = useState(false)
  
  // Create Form State
  const [name, setName] = useState('')
  const [partitions, setPartitions] = useState(3)
  const [errorMsg, setError] = useState('')

  const { data: topics = [], isLoading } = useQuery({
    queryKey: ['topics'],
    queryFn: () => api.get('/topics').then(res => res.data)
  })

  const createMutation = useMutation({
    mutationFn: (newTopic: any) => api.post('/topics', newTopic),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['topics'] })
      setShowCreate(false)
      setName('')
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to create topic')
    }
  })

  const approveMutation = useMutation({
    mutationFn: (topicName: string) => api.post(`/topics/${topicName}/approve`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['topics'] })
    }
  })

  const deleteMutation = useMutation({
    mutationFn: (topicName: string) => api.delete(`/topics/${topicName}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['topics'] })
    }
  })

  // Topics Filter
  const visibleTopics = topics.filter((t: any) => user?.role === 'ADMIN' || t.requestedBy === user?.id)

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    createMutation.mutate({
      name,
      cluster: 'local',
      partitions: Number(partitions),
      replicas: 1
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Topics</h2>
          <p className="text-muted-foreground">Manage topic metadata and request allocations.</p>
        </div>
        <Button onClick={() => setShowCreate(!showCreate)} className="gap-2">
          <Plus className="h-4 w-4" />
          {user?.role === 'ADMIN' ? 'Add Topic' : 'Request Topic'}
        </Button>
      </div>

      {showCreate && (
        <Card className="border-primary/50 shadow-md">
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <ShieldAlert className="h-5 w-5 text-primary" />
              New Topic Initialization
            </h3>
            <form onSubmit={handleCreate} className="space-y-4 max-w-xl">
              <div>
                <label className="text-sm font-medium mb-1 block">Topic Name</label>
                <Input value={name} onChange={e => setName(e.target.value)} placeholder="orders.events.v1" required />
              </div>
              <div>
                <label className="text-sm font-medium mb-1 block">Partitions</label>
                <Input type="number" min="1" value={partitions} onChange={e => setPartitions(Number(e.target.value))} required />
              </div>
              
              <div className="bg-secondary/50 p-3 rounded-md border text-sm text-muted-foreground flex gap-3">
                <CheckCircle2 className="h-4 w-4 mt-0.5 text-primary" />
                <div>
                  <p className="font-medium text-foreground">Cedar Evaluation Pre-Check</p>
                  <p>Action <code>CreateTopic</code> will be evaluated against your context ({user?.username}). Initial state will be PENDING.</p>
                </div>
              </div>

              {errorMsg && <p className="text-sm text-destructive">{errorMsg}</p>}
              <div className="pt-2 flex gap-2">
                <Button type="submit" disabled={createMutation.isPending}>
                  {createMutation.isPending ? 'Submitting...' : (user?.role === 'ADMIN' ? 'Add Topic' : 'Propose Topic')}
                </Button>
                <Button type="button" variant="outline" className="bg-transparent border-input" onClick={() => setShowCreate(false)}>Cancel</Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      <div className="border rounded-md border-border bg-card">
        <div className="overflow-x-auto">
          <table className="w-full text-sm text-left">
            <thead className="bg-secondary/50 text-muted-foreground">
              <tr>
                <th className="px-4 py-3 font-medium">Topic Name</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Partitions</th>
                <th className="px-4 py-3 font-medium">Owner</th>
                <th className="px-4 py-3 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td colSpan={5} className="p-4 text-center text-muted-foreground">Loading topics...</td></tr>
              ) : visibleTopics.map((topic: any) => (
                <tr key={topic.id} className="border-t border-border hover:bg-muted/50 transition-colors">
                  <td className="px-4 py-3 font-medium text-primary cursor-pointer hover:underline" onClick={() => window.open(`http://localhost:6060/ui/clusters/local/all-topics/${topic.name}`, '_blank')}>
                    {topic.name}
                  </td>
                  <td className="px-4 py-3">
                    <Badge className={topic.status === 'APPROVED' ? 'bg-primary/20 text-primary hover:bg-primary/30 border-transparent' : 'bg-secondary text-muted-foreground border-border'}>
                      {topic.status}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-muted-foreground">{topic.partitions} p | {topic.replicas} r</td>
                  <td className="px-4 py-3">{topic.requestedBy}</td>
                  <td className="px-4 py-3">
                    <div className="flex gap-2">
                      {topic.status === 'PENDING' && user?.role === 'ADMIN' && (
                        <Button 
                          onClick={() => approveMutation.mutate(topic.name)}
                          disabled={approveMutation.isPending}
                          className="h-7 px-3 py-0 text-xs bg-primary/10 text-primary hover:bg-primary/20"
                        >
                          Approve
                        </Button>
                      )}
                      
                      {user?.role === 'ADMIN' && (
                         <Button 
                           variant="outline"
                           onClick={() => {
                             // Minimal mock logic if API endpoint missing in backend
                             deleteMutation.mutate(topic.name)
                           }}
                           className="h-7 px-3 py-0 text-xs text-destructive border-destructive/50 hover:bg-destructive/10"
                         >
                           Delete
                         </Button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
              {visibleTopics.length === 0 && !isLoading && (
                <tr><td colSpan={5} className="p-8 text-center text-muted-foreground">No topics found for your access level.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
