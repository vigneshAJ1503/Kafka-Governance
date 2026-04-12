'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/core'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { ListTree, FileLock2, Clock } from 'lucide-react'
import Link from 'next/link'

export default function Dashboard() {
  const { data: topics = [] } = useQuery({
    queryKey: ['topics'],
    queryFn: () => api.get('/topics').then(res => res.data)
  })

  // Simulated metrics since we don't have a direct policies GET api
  const pendingApprovals = topics.filter((t: any) => t.status === 'PENDING').length

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Overview</h2>
        <p className="text-muted-foreground">Monitor cluster topic allocations and recent policy activity.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Link href="/topics" className="block transition-transform hover:scale-[1.02]">
          <Card className="h-full">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Total Topics</CardTitle>
              <ListTree className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{topics.length}</div>
              <p className="text-xs text-muted-foreground mt-1">Across 1 cluster</p>
            </CardContent>
          </Card>
        </Link>
        
        <Link href="/topics" className="block transition-transform hover:scale-[1.02]">
          <Card className="h-full border-primary/20 hover:border-primary/50 transition-colors">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Pending Approvals</CardTitle>
              <Clock className="h-4 w-4 text-primary" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-primary">{pendingApprovals}</div>
              <p className="text-xs text-muted-foreground mt-1">Require Admin Review</p>
            </CardContent>
          </Card>
        </Link>

        <Link href="/policies" className="block transition-transform hover:scale-[1.02]">
          <Card className="h-full">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Active Policies</CardTitle>
              <FileLock2 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">1</div>
              <p className="text-xs text-muted-foreground mt-1">Configured Cedar Rules</p>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Mock Activity Feed to represent density */}
      <h3 className="text-lg font-semibold tracking-tight mt-8">Recent Activity Loop</h3>
      <div className="border rounded-lg bg-card">
        {topics.slice(0, 5).map((topic: any, idx: number) => (
          <div key={idx} className="flex items-center gap-4 p-4 border-b last:border-0 border-border">
            <div className="w-2 h-2 rounded-full bg-primary" />
            <div className="flex-1">
              <p className="text-sm font-medium">Topic <span className="text-primary">{topic.name}</span> created by {topic.requestedBy}</p>
              <p className="text-xs text-muted-foreground">{new Date(topic.createdAt).toLocaleString()}</p>
            </div>
            <div className={`text-xs px-2 py-1 rounded border ${topic.status === 'APPROVED' ? 'bg-primary/20 text-primary border-primary/20' : 'bg-secondary text-muted-foreground'}`}>
              {topic.status}
            </div>
          </div>
        ))}
        {topics.length === 0 && <div className="p-8 text-center text-muted-foreground text-sm">No activity recorded yet</div>}
      </div>
    </div>
  )
}
