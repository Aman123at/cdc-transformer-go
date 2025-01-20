'use client'

import { useState } from 'react'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import DatabaseSection from './DatabaseSection'

export default function ResizableDrawers() {
  const [postgresStatus, setPostgresStatus] = useState('connected')
  const [mongodbStatus, setMongodbStatus] = useState('connected')

  return (
    <ResizablePanelGroup direction="horizontal" className="h-full">
      <ResizablePanel defaultSize={50}>
        <div className="p-4">
          <DatabaseSection
            title="PostgreSQL"
            status={postgresStatus}
            onStatusChange={setPostgresStatus}
            type="postgres"
          />
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle className="bg-secondary hover:bg-primary transition-colors" />
      <ResizablePanel defaultSize={50}>
        <div className="p-4">
          <DatabaseSection
            title="MongoDB"
            status={mongodbStatus}
            onStatusChange={setMongodbStatus}
            type="mongodb"
          />
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  )
}

