'use client'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'

export default function TransformLogsModal({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) {
  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Transform Logs</DialogTitle>
        </DialogHeader>
        <ScrollArea className="h-[300px] w-full rounded-md border p-4">
          {/* <pre className="text-sm">
            {`[2023-04-01 10:00:00] Transformation started
[2023-04-01 10:00:01] Processing table 'users'
[2023-04-01 10:00:02] 100 records transformed
[2023-04-01 10:00:03] Processing table 'orders'
[2023-04-01 10:00:04] 50 records transformed
[2023-04-01 10:00:05] Transformation completed`}
          </pre> */}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}

