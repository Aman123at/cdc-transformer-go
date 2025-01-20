'use client'
import ResizableDrawers from '@/components/custom/Drawer'
import Header from '@/components/custom/Header'
import TransformLogsModal from '@/components/custom/TransformLogsModal'
import { useEffect, useState } from 'react'

export default function Home() {
  const [isLogsModalOpen, setIsLogsModalOpen] = useState(false)

  return (
    <div className="flex flex-col h-screen">
      <Header onCheckLogs={() => setIsLogsModalOpen(true)} />
      <div className="flex-grow">
        <ResizableDrawers />
      </div>
      <TransformLogsModal isOpen={isLogsModalOpen} onClose={() => setIsLogsModalOpen(false)} />
    </div>
  )
}

