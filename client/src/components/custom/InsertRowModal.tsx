'use client'

import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { IColumn, IRow } from '@/interfaces/commonInterface'

type InsertRowModalProps = {
  isOpen: boolean
  onClose: () => void
  columns: IColumn[]
  onInsertRow: (row: IRow) => void
}

export default function InsertRowModal({ isOpen, onClose, columns, onInsertRow }: InsertRowModalProps) {
  const [rowData, setRowData] = useState<IRow>({})

  const handleInputChange = (columnName: string,columnType:string, value: string) => {
    let newval = columnType==="integer" ? parseInt(value) : value 
    setRowData({ ...rowData, [columnName]: newval })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onInsertRow(rowData)
    setRowData({})
  }

  const getDisabled = (column:IColumn):boolean=>{
    let disabled = false
    if (column.constraint){
      if(column.constraint.primary_key){
        disabled = true
      }
    }
    return disabled
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Insert Row</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          {columns.map((column) => (
            <div key={column.name} className="space-y-2">
              {!getDisabled(column) && 
              <>
              <Label htmlFor={column.name}>{column.name}</Label>
              <Input
                id={column.name}
                placeholder={`Enter ${column.name}`}
                value={rowData[column.name] || ''}
                onChange={(e) => handleInputChange(column.name,column.type, e.target.value)}
                disabled={getDisabled(column)}
              />
              </>
              }
            </div>
          ))}
          <Button type="submit" className="w-full">Insert</Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}

