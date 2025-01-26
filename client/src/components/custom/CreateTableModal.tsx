'use client'

import { useEffect, useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ScrollArea } from '@/components/ui/scroll-area'
import { X, Plus } from 'lucide-react'
import { createTable, getTables } from '@/apiCalls/commonCalls'
import {v4 as uuidv4} from "uuid"
import { usePG } from '@/contexts/pgcontext'

type Column = {
  name: string
  type: string
}

type CreateTableModalProps = {
  isOpen: boolean
  onClose: () => void
}

const dataTypes = ['int', 'text', 'varchar(255)', 'boolean',  'timestamp',  'double']

export default function CreateTableModal({ isOpen, onClose }: CreateTableModalProps) {
  const [tableName, setTableName] = useState('')
  const [disableCreateBtn,setDisbaleCreateBtn] = useState<boolean>(false)
  const {fetchTables} = usePG()
  const [columns, setColumns] = useState<Column[]>([
    { name: 'id', type: 'int' }
  ])

  const handleAddColumn = () => {
    if (columns.length < 5) {
      setColumns([...columns, { name: '', type: '' }])
    }
  }

  useEffect(()=>{
    if(columns.length<2){
      setDisbaleCreateBtn(true)
    }
  },[columns])

  const handleRemoveColumn = (index: number) => {
    setColumns(columns.filter((_, i) => i !== index))
  }

  const handleColumnChange = (index: number, field: 'name' | 'type', value: string) => {
    if (value.length>15){
      setDisbaleCreateBtn(true)
    }else{
      setDisbaleCreateBtn(false)
    }
    const newColumns = [...columns]
    newColumns[index][field] = value
    setColumns(newColumns)
  }

  const cleanUp =()=>{
    setTableName('')
      setColumns([{ name: 'id', type: 'int'}])
      onClose()
  }

  const handleSubmit = async(e: React.FormEvent) => {
    e.preventDefault()
    if (tableName.trim() && columns.every(col => col.name && col.type)) {
      let sessionid = localStorage.getItem("cdc-session-id")
      if (!sessionid) {
        sessionid = uuidv4()
        localStorage.setItem("cdc-session-id",sessionid)
      }
      const {data,err}:any = await createTable({
        name:tableName.trim(),
        columns,
        sessionid
      })
      if(err){
        // show toast
        cleanUp()
      }
      if(data){
        console.log("Table created successfully")
        cleanUp()
        setTimeout(() => {
          
          fetchTables()
        }, 1000);
      }
      
    }
  }

  const isAddColumnDisabled = columns.length >= 5

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Table</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            placeholder="Enter Table Name"
            value={tableName}
            onChange={(e) => setTableName(e.target.value)}
          />
          <ScrollArea className="h-[200px] w-full rounded-md border p-4">
            {columns.map((column, index) => (
              <div key={index} className="flex items-center space-x-2 mb-2">
                <Input
                  placeholder="Column name"
                  value={column.name}
                  onChange={(e) => handleColumnChange(index, 'name', e.target.value)}
                  disabled={index === 0}
                />
                <Select
                  value={column.type}
                  onValueChange={(value) => handleColumnChange(index, 'type', value)}
                  disabled={index === 0}
                >
                  <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent>
                    {dataTypes.map((type) => (
                      <SelectItem key={type} value={type}>
                        {type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {index !== 0 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => handleRemoveColumn(index)}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            ))}
          </ScrollArea>
          <Button
            type="button"
            onClick={handleAddColumn}
            disabled={isAddColumnDisabled}
            className="w-full"
          >
            <Plus className="mr-2 h-4 w-4" /> Add Column
          </Button>
          <Button type="submit" className="w-full" disabled={disableCreateBtn}>Create Table</Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}

