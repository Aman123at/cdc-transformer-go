'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import TableAccordion from './TableAccordion'
import CollectionAccordion from './CollectionAccordion'
import CreateTableModal from './CreateTableModal'
import { PGContextType, usePG } from '@/contexts/pgcontext'
import { ICollection, ITable } from '@/interfaces/commonInterface'
import { MongoContextType, useMongo } from '@/contexts/mongocontext'

type DatabaseSectionProps = {
  title: string
  status: string
  onStatusChange: (status: string) => void
  type: 'postgres' | 'mongodb'
}

export default function DatabaseSection({ title, status, onStatusChange, type }: DatabaseSectionProps) {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const {tables}:PGContextType = usePG()
  const {collections}:MongoContextType = useMongo()
  
  return (
    <div className="h-full flex flex-col">
      <div className="flex justify-between items-center p-4 bg-secondary">
        <h2 className="text-lg font-semibold">{title}</h2>
        <div className="flex items-center">
          <span className={`mr-2 px-2 py-1 rounded-full text-xs ${
            status === 'connected' ? 'bg-green-500' : 
            status === 'disconnected' ? 'bg-red-500' : 'bg-yellow-500'
          }`}>
            {status}
          </span>
          {type === 'postgres' && (
            <Button size="sm" onClick={() => setIsCreateModalOpen(true)}>
              <Plus className="mr-1 h-4 w-4" /> Create Table
            </Button>
          )}
        </div>
      </div>
      <div className="flex-grow overflow-auto">
        <div className="p-2">
          {type === 'postgres' ? (
            tables ? tables.map((table:ITable, index:number) => (
              <TableAccordion key={index} tableName={table.tablename} columns={table.columns} rows={table.rows} />
            )):<></>
          ) : (
            collections ? collections.map((collection:ICollection, index:number) => (
              <CollectionAccordion key={index} collectionName={collection.collectionname} documents={collection.documents} />
            )):<></>
          )}
        </div>
      </div>
      {type === 'postgres' && (
        <CreateTableModal
          isOpen={isCreateModalOpen}
          onClose={() => setIsCreateModalOpen(false)}
        />
      )}
    </div>
  )
}

