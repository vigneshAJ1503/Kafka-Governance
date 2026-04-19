'use client'

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/core'
import { HardDrive } from 'lucide-react'

export default function AuditLogsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Audit Logs</h2>
        <p className="text-muted-foreground">Monitor system-wide events, permissions history, and access audits.</p>
      </div>

      <Card className="border-dashed border-2 bg-background/50">
        <CardContent className="flex flex-col items-center justify-center h-64 text-center">
          <HardDrive className="h-12 w-12 text-muted-foreground mb-4 opacity-50" />
          <h3 className="text-lg font-semibold text-foreground">Under Development</h3>
          <p className="text-muted-foreground text-sm max-w-sm mt-2">
            The Audit Logs ingestion module is currently being integrated with open-source streaming platforms (e.g. OpenSearch / ELK). 
            Check back later for historic tracking!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
