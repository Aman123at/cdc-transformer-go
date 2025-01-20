'use client'

import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

type Column = {
  name: string
  type: string
}

type EditColumnModalProps = {
  isOpen: boolean
  onClose: () => void
  column: Column
  onUpdateColumn: (updatedColumn: Column) => void
}

export default function EditColumnModal({ isOpen, onClose, column, onUpdateColumn }: EditColumnModalProps) {
  const [editedColumn, setEditedColumn] = useState(column)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onUpdateColumn(editedColumn)
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Column</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            placeholder="Column Name"
            value={editedColumn.name}
            onChange={(e) => setEditedColumn({ ...editedColumn, name: e.target.value })}
          />
          <Input
            placeholder="Column Type"
            value={editedColumn.type}
            onChange={(e) => setEditedColumn({ ...editedColumn, type: e.target.value })}
          />
          <Button type="submit">Update</Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}

